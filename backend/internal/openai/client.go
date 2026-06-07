package openai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const defaultChatModel = "gpt-4o-mini"
const defaultImageModel = "gpt-image-1"
const defaultImageQuality = "medium"

// Client calls OpenAI chat and image APIs.
type Client struct {
	apiKey       string
	chatModel    string
	imageModel   string
	imageQuality string
	httpClient   *http.Client
}

// NewClientFromEnv returns a client when OPENAI_API_KEY is set.
func NewClientFromEnv() (*Client, bool) {
	key := strings.TrimSpace(os.Getenv("OPENAI_API_KEY"))
	if key == "" {
		return nil, false
	}
	chatModel := strings.TrimSpace(os.Getenv("OPENAI_CHAT_MODEL"))
	if chatModel == "" {
		chatModel = defaultChatModel
	}
	imageModel := strings.TrimSpace(os.Getenv("OPENAI_IMAGE_MODEL"))
	if imageModel == "" {
		imageModel = defaultImageModel
	}
	imageQuality := strings.TrimSpace(os.Getenv("OPENAI_IMAGE_QUALITY"))
	if imageQuality == "" {
		imageQuality = defaultImageQuality
	}
	return &Client{
		apiKey:       key,
		chatModel:    chatModel,
		imageModel:   imageModel,
		imageQuality: imageQuality,
		httpClient:   &http.Client{Timeout: 120 * time.Second},
	}, true
}

type chatRequest struct {
	Model          string          `json:"model"`
	Messages       []chatMessage   `json:"messages"`
	ResponseFormat *responseFormat `json:"response_format,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type chatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *apiError `json:"error"`
}

type apiError struct {
	Message string `json:"message"`
}

// ChatJSON sends a prompt and returns parsed JSON bytes from the model response.
func (c *Client) ChatJSON(ctx context.Context, systemPrompt, userPrompt string) (json.RawMessage, error) {
	body, err := json.Marshal(chatRequest{
		Model: c.chatModel,
		Messages: []chatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		ResponseFormat: &responseFormat{Type: "json_object"},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("openai chat: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 400))
	}

	var parsed chatResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("openai chat: %s", parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return nil, fmt.Errorf("openai chat: empty response")
	}
	content := strings.TrimSpace(parsed.Choices[0].Message.Content)
	if content == "" {
		return nil, fmt.Errorf("openai chat: empty content")
	}
	var raw json.RawMessage
	if err := json.Unmarshal([]byte(content), &raw); err != nil {
		return nil, fmt.Errorf("openai chat: invalid JSON: %w", err)
	}
	return raw, nil
}

type imageRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	N       int    `json:"n"`
	Size    string `json:"size"`
	Quality string `json:"quality,omitempty"`
}

type imageResponse struct {
	Data []struct {
		URL     string `json:"url"`
		B64JSON string `json:"b64_json"`
	} `json:"data"`
	Error *apiError `json:"error"`
}

// GenerateImage creates an image and returns PNG bytes.
// GPT Image models return base64-encoded data; legacy DALL-E models may return URLs.
func (c *Client) GenerateImage(ctx context.Context, prompt string) ([]byte, error) {
	body, err := json.Marshal(imageRequest{
		Model:   c.imageModel,
		Prompt:  prompt,
		N:       1,
		Size:    "1024x1024",
		Quality: c.imageQuality,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/images/generations", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("openai image: HTTP %d: %s", resp.StatusCode, truncate(string(respBody), 400))
	}

	var parsed imageResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return nil, err
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("openai image: %s", parsed.Error.Message)
	}
	if len(parsed.Data) == 0 {
		return nil, fmt.Errorf("openai image: empty response")
	}
	if b64 := strings.TrimSpace(parsed.Data[0].B64JSON); b64 != "" {
		return base64.StdEncoding.DecodeString(b64)
	}
	if url := strings.TrimSpace(parsed.Data[0].URL); url != "" {
		return DownloadImage(ctx, url)
	}
	return nil, fmt.Errorf("openai image: empty response")
}

// DownloadImage fetches image bytes from a URL returned by OpenAI.
func DownloadImage(ctx context.Context, rawURL string) ([]byte, error) {
	if !allowedImageDownloadURL(rawURL) {
		return nil, fmt.Errorf("download image: disallowed URL host")
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("download image: HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(io.LimitReader(resp.Body, 8<<20))
}

func allowedImageDownloadURL(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	allowedHosts := []string{
		"oaidalleapiprodscus.blob.core.windows.net",
		"openai.com",
	}
	for _, allowed := range allowedHosts {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "…"
}
