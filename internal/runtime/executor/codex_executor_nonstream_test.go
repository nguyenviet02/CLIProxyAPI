package executor

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/router-for-me/CLIProxyAPI/v6/internal/config"
	cliproxyauth "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/auth"
	cliproxyexecutor "github.com/router-for-me/CLIProxyAPI/v6/sdk/cliproxy/executor"
)

func TestCodexExecutorExecute_AcceptsDirectResponseObject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/responses" {
			t.Fatalf("path = %s, want /responses", r.URL.Path)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if len(body) == 0 {
			t.Fatal("expected request body")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"resp_direct","object":"response","status":"completed","output":[{"type":"message","id":"msg_1","content":[{"type":"output_text","text":"hello from file"}]}],"usage":{"input_tokens":10,"output_tokens":5,"total_tokens":15}}`))
	}))
	defer server.Close()

	executor := NewCodexExecutor(&config.Config{})
	auth := &cliproxyauth.Auth{
		Provider: "codex",
		Attributes: map[string]string{
			"api_key":  "test-key",
			"base_url": server.URL,
		},
	}
	req := cliproxyexecutor.Request{
		Model:   "gpt-5.3-codex",
		Payload: []byte(`{"model":"gpt-5.3-codex","input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hello"}]}]}`),
	}

	resp, err := executor.Execute(context.Background(), auth, req, cliproxyexecutor.Options{SourceFormat: "openai-response"})
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if string(resp.Payload) == "" {
		t.Fatal("expected translated payload")
	}
	if got := string(resp.Payload); got == "[]" {
		t.Fatalf("unexpected empty payload: %s", got)
	}
}
