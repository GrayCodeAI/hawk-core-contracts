// Package sessions defines cross-repo types for agent session state management.
//
// Every multi-step coding-agent pipeline (localize → repair → validate) needs a
// shared vocabulary for session identity, pipeline phase, context snapshots, and
// cost accounting. Defining these types here — in the leaf contracts module — lets
// all engines and the orchestrator agree on a single schema without creating
// circular dependencies.
//
// Evidence basis: CAT (arXiv 2512.22087) demonstrates that context management as
// a first-class concern improves long-horizon agent performance. ACON
// (arXiv 2510.00615) shows 26–54% token savings when phase-aware compression is
// applied. The Tokenomics study (arXiv 2601.14470) measured that code review
// alone consumes 59.4% of total tokens in multi-agent SE systems.
package sessions

import (
	"fmt"
	"time"
)

// Phase identifies the pipeline stage that generated a token-usage event.
// Engines tag their tok.Tracker events with the appropriate Phase so that
// hawk's orchestrator can build per-phase cost breakdowns and apply optimal
// compaction budgets for each stage.
type Phase string

const (
	// PhaseLocalize is the file/symbol retrieval stage (Agentless Stage 1+2).
	PhaseLocalize Phase = "localize"
	// PhaseRepair is the patch-generation stage (Agentless Stage 2).
	PhaseRepair Phase = "repair"
	// PhaseValidate is the test-execution and patch-ranking stage (Agentless Stage 3).
	PhaseValidate Phase = "validate"
	// PhaseReview is the LLM-based code-review stage (sight/inspect).
	// The Tokenomics paper shows this phase dominates token usage (~59.4%).
	PhaseReview Phase = "review"
	// PhasePlanning is the task-decomposition stage (pre-execution).
	PhasePlanning Phase = "planning"
	// PhaseUnknown is the default zero value — used when phase attribution is unavailable.
	PhaseUnknown Phase = ""
)

// ParsePhase parses a phase name string into a Phase constant.
// Returns PhaseUnknown for unrecognised values rather than an error, so
// callers that receive phase names from JSON/TOML do not need to handle errors.
func ParsePhase(s string) Phase {
	switch Phase(s) {
	case PhaseLocalize, PhaseRepair, PhaseValidate, PhaseReview, PhasePlanning:
		return Phase(s)
	default:
		return PhaseUnknown
	}
}

// String returns the string representation of a Phase.
func (p Phase) String() string {
	if p == PhaseUnknown {
		return "unknown"
	}
	return string(p)
}

// SessionID is an opaque, globally unique identifier for a single agent session.
type SessionID string

// ContextSnapshot is a point-in-time, compressed view of an agent's working
// memory. It is created by the tok CompactContext call and stored so that a
// new session can resume from a known good state without replaying history.
//
// Anchors hold immutable facts (file paths, error messages, invariants) that
// must survive every compaction and must always be re-injected into context.
type ContextSnapshot struct {
	SessionID    SessionID `json:"session_id"`
	Phase        Phase     `json:"phase"`
	TokenCount   int       `json:"token_count"`
	CompressedAt time.Time `json:"compressed_at"`
	Summary      string    `json:"summary,omitempty"`
	Anchors      []string  `json:"anchors,omitempty"`
}

// ToolCallRecord captures a single tool invocation within a session, tagged
// with the pipeline phase in which it occurred. Accumulating these records
// lets hawk build a per-tool, per-phase cost profile that drives adaptive
// budget allocation in subsequent runs.
type ToolCallRecord struct {
	SessionID    SessionID `json:"session_id"`
	Phase        Phase     `json:"phase"`
	ToolName     string    `json:"tool_name"`
	InputTokens  int       `json:"input_tokens"`
	OutputTokens int       `json:"output_tokens"`
	DurationMs   int64     `json:"duration_ms"`
	Error        string    `json:"error,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// PhaseUsage captures token and cost usage for a single pipeline phase.
type PhaseUsage struct {
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	TotalTokens  int     `json:"total_tokens"`
	CostUSD      float64 `json:"cost_usd"`
}

// CostAccumulator tracks token consumption and monetary cost across all
// pipeline phases within a session. It is the single source of truth for
// the cost-per-resolution metric exposed by hawk's ResolutionPipeline.
type CostAccumulator struct {
	SessionID    SessionID            `json:"session_id"`
	ByPhase      map[Phase]PhaseUsage `json:"by_phase"`
	TotalTokens  int                  `json:"total_tokens"`
	TotalCostUSD float64              `json:"total_cost_usd"`
}

// NewCostAccumulator returns a zero-value CostAccumulator for the given session.
func NewCostAccumulator(id SessionID) *CostAccumulator {
	return &CostAccumulator{
		SessionID: id,
		ByPhase:   make(map[Phase]PhaseUsage),
	}
}

// Add records input/output token usage and cost for a given phase.
func (c *CostAccumulator) Add(phase Phase, inputTokens, outputTokens int, costUSD float64) {
	if c.ByPhase == nil {
		c.ByPhase = make(map[Phase]PhaseUsage)
	}
	pu := c.ByPhase[phase]
	pu.InputTokens += inputTokens
	pu.OutputTokens += outputTokens
	pu.TotalTokens += inputTokens + outputTokens
	pu.CostUSD += costUSD
	c.ByPhase[phase] = pu
	c.TotalTokens += inputTokens + outputTokens
	c.TotalCostUSD += costUSD
}

// PhaseShare returns the fraction of total tokens attributed to the given phase.
// Returns 0 when no tokens have been recorded.
func (c *CostAccumulator) PhaseShare(phase Phase) float64 {
	if c.TotalTokens == 0 {
		return 0
	}
	pu, ok := c.ByPhase[phase]
	if !ok {
		return 0
	}
	return float64(pu.TotalTokens) / float64(c.TotalTokens)
}

// String returns a one-line summary suitable for log output.
func (c *CostAccumulator) String() string {
	return fmt.Sprintf("session=%s tokens=%d cost=$%.4f", c.SessionID, c.TotalTokens, c.TotalCostUSD)
}
