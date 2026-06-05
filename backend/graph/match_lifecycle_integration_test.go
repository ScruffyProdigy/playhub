package graph

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/scruffyprodigy/playhub/internal/auth"
	"github.com/scruffyprodigy/playhub/internal/store"
)

func TestReturnDestinationAfterMatchmaking(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoQueue(t, env.Store)

	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")
	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "lifecycle-pepper")

	provisioner := &syncProvisioner{}
	env.resolverWithProvisioner(t, provisioner)

	_, cookieA := createTestUserSession(t, ctx, env, cleaner)
	_, cookieB := createTestUserSession(t, ctx, env, cleaner)

	joinQuery := `mutation Join($id: ID!) { joinQueue(queueId: $id) { queued queuedCount } }`
	vars := map[string]any{"id": demoDefaultQueueID}
	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	call := provisioner.lastCall()
	matchID := call.Assignment.ExternalMatchID
	if matchID == "" {
		t.Fatal("expected external match id from provision")
	}

	destQuery := `query Return($matchId: ID!) {
		returnDestination(matchId: $matchId) { path kind }
	}`
	body := postGraphQL(t, env.Handler, destQuery, map[string]any{"matchId": matchID}, cookieA)

	var resp struct {
		Data struct {
			ReturnDestination struct {
				Path string `json:"path"`
				Kind string `json:"kind"`
			} `json:"returnDestination"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode: %v body=%s", err, body)
	}
	if resp.Data.ReturnDestination.Path != "/" {
		t.Fatalf("path = %q, want /", resp.Data.ReturnDestination.Path)
	}
	if resp.Data.ReturnDestination.Kind != store.ReturnKindCatalogLFG {
		t.Fatalf("kind = %q", resp.Data.ReturnDestination.Kind)
	}
}

func TestReportMatchResultCompletesSession(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoQueue(t, env.Store)

	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")
	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "lifecycle-pepper")

	provisioner := &syncProvisioner{}
	env.resolverWithProvisioner(t, provisioner)

	_, cookieA := createTestUserSession(t, ctx, env, cleaner)
	_, cookieB := createTestUserSession(t, ctx, env, cleaner)

	joinQuery := `mutation Join($id: ID!) { joinQueue(queueId: $id) { queued queuedCount } }`
	vars := map[string]any{"id": demoDefaultQueueID}
	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	matchID := provisioner.lastCall().Assignment.ExternalMatchID
	sessionID, err := uuid.Parse(matchID)
	if err != nil {
		t.Fatalf("parse session id: %v", err)
	}

	gameID := uuid.MustParse(store.DemoPrimaryGameIDStr)
	serviceToken, err := auth.FormatGameServiceToken(gameID)
	if err != nil {
		t.Fatalf("FormatGameServiceToken: %v", err)
	}

	mutation := `mutation Report($matchId: ID!, $status: MatchResultStatus!) {
		reportMatchResult(matchId: $matchId, status: $status)
	}`
	reportBody := postGraphQLWithBearer(t, env.Handler, serviceToken, mutation, map[string]any{
		"matchId": matchID,
		"status":  "COMPLETED",
	})
	var reportResp struct {
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
		Data struct {
			ReportMatchResult bool `json:"reportMatchResult"`
		} `json:"data"`
	}
	if err := json.Unmarshal(reportBody, &reportResp); err != nil {
		t.Fatalf("decode report: %v body=%s", err, reportBody)
	}
	if len(reportResp.Errors) > 0 {
		t.Fatalf("reportMatchResult errors: %+v", reportResp.Errors)
	}
	if !reportResp.Data.ReportMatchResult {
		t.Fatalf("reportMatchResult = false, body=%s", reportBody)
	}

	session, err := env.Store.GetSessionByID(ctx, sessionID)
	if err != nil {
		t.Fatalf("GetSessionByID: %v", err)
	}
	if session.Status != "completed" {
		t.Fatalf("session status = %q, want completed", session.Status)
	}
}

func TestFullGameLoopReturnAndRequeue(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoQueue(t, env.Store)

	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")
	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "lifecycle-pepper")

	provisioner := &syncProvisioner{}
	env.resolverWithProvisioner(t, provisioner)

	_, cookieA := createTestUserSession(t, ctx, env, cleaner)
	_, cookieB := createTestUserSession(t, ctx, env, cleaner)

	joinQuery := `mutation Join($id: ID!) { joinQueue(queueId: $id) { queued joinUrl sessionId } }`
	vars := map[string]any{"id": demoDefaultQueueID}

	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	matchBody := postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	var matchResp struct {
		Data struct {
			JoinQueue struct {
				SessionID *string `json:"sessionId"`
				JoinURL   *string `json:"joinUrl"`
			} `json:"joinQueue"`
		} `json:"data"`
	}
	if err := json.Unmarshal(matchBody, &matchResp); err != nil {
		t.Fatalf("decode match join: %v body=%s", err, matchBody)
	}
	if matchResp.Data.JoinQueue.SessionID == nil {
		t.Fatalf("expected sessionId on match, got %s", matchBody)
	}
	firstSessionID := *matchResp.Data.JoinQueue.SessionID

	activeQuery := `query { myActiveQueue { queueId status sessionId joinUrl } }`
	bodyBefore := postGraphQL(t, env.Handler, activeQuery, nil, cookieA)
	var beforeResp struct {
		Data struct {
			MyActiveQueue *struct {
				Status    string  `json:"status"`
				SessionID *string `json:"sessionId"`
			} `json:"myActiveQueue"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyBefore, &beforeResp); err != nil {
		t.Fatalf("decode active before: %v", err)
	}
	if beforeResp.Data.MyActiveQueue == nil || beforeResp.Data.MyActiveQueue.Status != "MATCHED" {
		t.Fatalf("expected MATCHED before return, got %s", bodyBefore)
	}
	if beforeResp.Data.MyActiveQueue.SessionID == nil || *beforeResp.Data.MyActiveQueue.SessionID != firstSessionID {
		t.Fatalf("active session = %v, want %s", beforeResp.Data.MyActiveQueue.SessionID, firstSessionID)
	}

	destQuery := `query Return($matchId: ID!) { returnDestination(matchId: $matchId) { path kind } }`
	postGraphQL(t, env.Handler, destQuery, map[string]any{"matchId": firstSessionID}, cookieA)
	postGraphQL(t, env.Handler, destQuery, map[string]any{"matchId": firstSessionID}, cookieB)

	bodyAfterReturn := postGraphQL(t, env.Handler, activeQuery, nil, cookieA)
	var afterReturnResp struct {
		Data struct {
			MyActiveQueue *struct {
				Status string `json:"status"`
			} `json:"myActiveQueue"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyAfterReturn, &afterReturnResp); err != nil {
		t.Fatalf("decode active after return: %v", err)
	}
	if afterReturnResp.Data.MyActiveQueue != nil {
		t.Fatalf("expected no active queue after return, got %s", bodyAfterReturn)
	}

	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	requeueBody := postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	var requeueResp struct {
		Data struct {
			JoinQueue struct {
				SessionID *string `json:"sessionId"`
			} `json:"joinQueue"`
		} `json:"data"`
	}
	if err := json.Unmarshal(requeueBody, &requeueResp); err != nil {
		t.Fatalf("decode requeue join: %v body=%s", err, requeueBody)
	}
	if requeueResp.Data.JoinQueue.SessionID == nil {
		t.Fatalf("expected sessionId on requeue match, got %s", requeueBody)
	}
	secondSessionID := *requeueResp.Data.JoinQueue.SessionID
	if secondSessionID == firstSessionID {
		t.Fatalf("re-queue reused session %s", firstSessionID)
	}

	firstUUID, err := uuid.Parse(firstSessionID)
	if err != nil {
		t.Fatalf("parse first session: %v", err)
	}
	firstSession, err := env.Store.GetSessionByID(ctx, firstUUID)
	if err != nil {
		t.Fatalf("GetSessionByID first: %v", err)
	}
	if firstSession.Status != "completed" {
		t.Fatalf("first session status = %q, want completed", firstSession.Status)
	}

	bodyRequeued := postGraphQL(t, env.Handler, activeQuery, nil, cookieA)
	var requeuedActive struct {
		Data struct {
			MyActiveQueue *struct {
				Status    string  `json:"status"`
				SessionID *string `json:"sessionId"`
			} `json:"myActiveQueue"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyRequeued, &requeuedActive); err != nil {
		t.Fatalf("decode active after requeue: %v", err)
	}
	if requeuedActive.Data.MyActiveQueue == nil || requeuedActive.Data.MyActiveQueue.Status != "MATCHED" {
		t.Fatalf("expected MATCHED after requeue, got %s", bodyRequeued)
	}
	if requeuedActive.Data.MyActiveQueue.SessionID == nil || *requeuedActive.Data.MyActiveQueue.SessionID != secondSessionID {
		t.Fatalf("active session after requeue = %v, want %s", requeuedActive.Data.MyActiveQueue.SessionID, secondSessionID)
	}
}

func TestReturnDestinationClearsMatchedQueue(t *testing.T) {
	env := newQueueIntegrationEnv(t)
	cleaner := env.newCleaner(t)
	ctx := context.Background()
	clearDemoQueue(t, env.Store)

	t.Setenv("LOBBY_ISSUER_URL", "http://localhost:8080")
	t.Setenv("LOBBY_PUBLIC_URL", "http://localhost:5173")
	t.Setenv("LOBBY_GAME_TOKEN_PEPPER", "lifecycle-pepper")

	provisioner := &syncProvisioner{}
	env.resolverWithProvisioner(t, provisioner)

	_, cookieA := createTestUserSession(t, ctx, env, cleaner)
	_, cookieB := createTestUserSession(t, ctx, env, cleaner)

	joinQuery := `mutation Join($id: ID!) { joinQueue(queueId: $id) { queued queuedCount } }`
	vars := map[string]any{"id": demoDefaultQueueID}
	postGraphQL(t, env.Handler, joinQuery, vars, cookieA)
	postGraphQL(t, env.Handler, joinQuery, vars, cookieB)

	matchID := provisioner.lastCall().Assignment.ExternalMatchID

	activeBefore := `query { myActiveQueue { queueId status joinUrl } }`
	bodyBefore := postGraphQL(t, env.Handler, activeBefore, nil, cookieA)
	var beforeResp struct {
		Data struct {
			MyActiveQueue *struct {
				Status string `json:"status"`
			} `json:"myActiveQueue"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyBefore, &beforeResp); err != nil {
		t.Fatalf("decode before: %v", err)
	}
	if beforeResp.Data.MyActiveQueue == nil || beforeResp.Data.MyActiveQueue.Status != "MATCHED" {
		t.Fatalf("expected MATCHED before return, got %s", bodyBefore)
	}

	destQuery := `query Return($matchId: ID!) {
		returnDestination(matchId: $matchId) { path kind }
	}`
	postGraphQL(t, env.Handler, destQuery, map[string]any{"matchId": matchID}, cookieA)

	bodyAfter := postGraphQL(t, env.Handler, activeBefore, nil, cookieA)
	var afterResp struct {
		Data struct {
			MyActiveQueue *struct {
				Status string `json:"status"`
			} `json:"myActiveQueue"`
		} `json:"data"`
	}
	if err := json.Unmarshal(bodyAfter, &afterResp); err != nil {
		t.Fatalf("decode after: %v", err)
	}
	if afterResp.Data.MyActiveQueue != nil {
		t.Fatalf("expected no active queue after return hub, got %s", bodyAfter)
	}
}
