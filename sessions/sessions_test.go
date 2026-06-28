package sessions_test

import (
	"testing"

	"github.com/GrayCodeAI/hawk-core-contracts/sessions"
)

func TestParsePhase(t *testing.T) {
	cases := []struct {
		input string
		want  sessions.Phase
	}{
		{"localize", sessions.PhaseLocalize},
		{"repair", sessions.PhaseRepair},
		{"validate", sessions.PhaseValidate},
		{"review", sessions.PhaseReview},
		{"planning", sessions.PhasePlanning},
		{"", sessions.PhaseUnknown},
		{"bogus", sessions.PhaseUnknown},
	}
	for _, tc := range cases {
		got := sessions.ParsePhase(tc.input)
		if got != tc.want {
			t.Errorf("ParsePhase(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestPhaseString(t *testing.T) {
	if sessions.PhaseLocalize.String() != "localize" {
		t.Errorf("PhaseLocalize.String() = %q, want %q", sessions.PhaseLocalize.String(), "localize")
	}
	if sessions.PhaseUnknown.String() != "unknown" {
		t.Errorf("PhaseUnknown.String() = %q, want %q", sessions.PhaseUnknown.String(), "unknown")
	}
}

func TestCostAccumulator(t *testing.T) {
	acc := sessions.NewCostAccumulator("sess-1")

	acc.Add(sessions.PhaseReview, 5000, 200, 0.015)
	acc.Add(sessions.PhaseLocalize, 1000, 100, 0.003)
	acc.Add(sessions.PhaseRepair, 2000, 300, 0.006)

	total := 5000 + 200 + 1000 + 100 + 2000 + 300
	if acc.TotalTokens != total {
		t.Errorf("TotalTokens = %d, want %d", acc.TotalTokens, total)
	}

	// Review should be the dominant phase (~59% in real SE workloads)
	reviewShare := acc.PhaseShare(sessions.PhaseReview)
	if reviewShare <= 0 || reviewShare > 1 {
		t.Errorf("PhaseShare(review) = %f, want (0,1]", reviewShare)
	}

	// Zero-token accumulator returns 0 share
	empty := sessions.NewCostAccumulator("sess-empty")
	if empty.PhaseShare(sessions.PhaseReview) != 0 {
		t.Error("empty accumulator should return 0 share")
	}
}

func TestCostAccumulatorString(t *testing.T) {
	acc := sessions.NewCostAccumulator("test-session")
	acc.Add(sessions.PhaseReview, 100, 50, 0.001)
	s := acc.String()
	if s == "" {
		t.Error("String() must not be empty")
	}
}
