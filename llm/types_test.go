package llm

import (
	"context"
	"testing"

	"github.com/GrayCodeAI/hawk-core-contracts/tools"
)

func TestToolCallAliasesToolsPackage(t *testing.T) {
	// Ensure llm.ToolCall and tools.ToolCall are the same type (not merely similar).
	call := tools.ToolCall{
		ID:   "t1",
		Name: "read",
		Arguments: map[string]interface{}{
			"path": "main.go",
		},
	}
	var asLLM ToolCall = call
	if asLLM.ID != "t1" || asLLM.Name != "read" {
		t.Fatalf("ToolCall alias lost fields: %+v", asLLM)
	}
	result := tools.ToolResult{
		ToolUseID: "t1",
		Content:   "ok",
		IsError:   false,
	}
	var asLLMResult ToolResult = result
	if asLLMResult.ToolUseID != "t1" || asLLMResult.Content != "ok" {
		t.Fatalf("ToolResult alias lost fields: %+v", asLLMResult)
	}
}

func TestNewStreamResultAndClose(t *testing.T) {
	events := make(chan EyrieStreamEvent, 1)
	events <- EyrieStreamEvent{Type: "content", Content: "hi"}
	close(events)

	cancelled := false
	sr := NewStreamResult(events, "req-1", func() { cancelled = true })
	if sr.RequestID != "req-1" {
		t.Fatalf("RequestID = %q, want req-1", sr.RequestID)
	}
	got := <-sr.Events
	if got.Content != "hi" {
		t.Fatalf("event content = %q, want hi", got.Content)
	}
	sr.Close()
	if !cancelled {
		t.Fatal("Close did not invoke cancel")
	}
	// Close is safe when cancel is nil.
	NewStreamResult(nil, "", nil).Close()
}

func TestStreamResultCloseNilReceiver(t *testing.T) {
	var sr *StreamResult
	sr.Close() // must not panic
}

func TestIntentAndModelClassConstants(t *testing.T) {
	if IntentFast == "" || IntentBalanced == "" || IntentReasoning == "" || IntentEconomical == "" {
		t.Fatal("Intent constants must be non-empty")
	}
	if ModelClassEconomical == "" || ModelClassBalanced == "" || ModelClassPremium == "" {
		t.Fatal("ModelClass constants must be non-empty")
	}
	if CheckOK != "ok" || CheckFail != "fail" || CheckWarn != "warn" {
		t.Fatalf("check statuses: ok=%q fail=%q warn=%q", CheckOK, CheckFail, CheckWarn)
	}
}

func TestEventStreamerInterfaceShape(t *testing.T) {
	// Compile-time: a stub that matches EventStreamer keeps the host port stable.
	var _ EventStreamer = stubStreamer{}
}

type stubStreamer struct{}

func (stubStreamer) Next() bool              { return false }
func (stubStreamer) Event() EyrieStreamEvent { return EyrieStreamEvent{} }
func (stubStreamer) Err() error              { return nil }
func (stubStreamer) Close() error            { return nil }

func TestGenerateRequestZeroValue(t *testing.T) {
	var req GenerateRequest
	if req.Messages != nil || req.SystemPrompt != "" {
		t.Fatalf("unexpected zero GenerateRequest: %+v", req)
	}
	_ = context.Background()
}
