package spiritanimal

import "testing"

func TestIsRetryableImageError(t *testing.T) {
	cases := []struct {
		err  string
		want bool
	}{
		{"openai image: HTTP 500: server_error", true},
		{"openai image: HTTP 503: upstream connect error", true},
		{"context deadline exceeded", true},
		{"openai image: HTTP 400: invalid prompt", false},
		{"openai image: HTTP 429: quota", false},
	}
	for _, tc := range cases {
		got := isRetryableImageError(errString(tc.err))
		if got != tc.want {
			t.Fatalf("isRetryableImageError(%q) = %v, want %v", tc.err, got, tc.want)
		}
	}
}

type errString string

func (e errString) Error() string { return string(e) }
