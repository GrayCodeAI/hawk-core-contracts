package agent

import (
	"strings"
	"testing"
)

func TestParseSubagentTypeAliases(t *testing.T) {
	cases := []struct {
		in   string
		want SubagentType
	}{
		{"", TypeExplore},
		{"explore", TypeExplore},
		{"plan", TypePlan},
		{"general", TypeGeneralPurpose},
		{"general-purpose", TypeGeneralPurpose},
		{"General_Purpose", TypeGeneralPurpose},
	}
	for _, tc := range cases {
		got, err := ParseSubagentType(tc.in)
		if err != nil {
			t.Fatalf("ParseSubagentType(%q): %v", tc.in, err)
		}
		if got != tc.want {
			t.Errorf("ParseSubagentType(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
	if _, err := ParseSubagentType("wizard"); err == nil {
		t.Fatal("expected error for unknown type")
	}
}

func TestParseCapabilityAndIsolation(t *testing.T) {
	cap, err := ParseCapabilityMode("ReadOnly")
	if err != nil || cap != CapReadOnly {
		t.Fatalf("cap: %v %v", cap, err)
	}
	cap, err = ParseCapabilityMode("all")
	if err != nil || cap != CapAll {
		t.Fatalf("cap all: %v %v", cap, err)
	}
	iso, err := ParseIsolationMode("")
	if err != nil || iso != IsoNone {
		t.Fatalf("iso empty: %v %v", iso, err)
	}
	iso, err = ParseIsolationMode("worktree")
	if err != nil || iso != IsoWorktree {
		t.Fatalf("iso wt: %v %v", iso, err)
	}
	if _, err := ParseIsolationMode("container"); err == nil {
		t.Fatal("expected isolation error")
	}
}

func TestNormalizeDefaults(t *testing.T) {
	n, err := (SpawnRequest{Prompt: "find TODOs"}).Normalize()
	if err != nil {
		t.Fatal(err)
	}
	if n.SubagentType != TypeExplore {
		t.Errorf("type=%q", n.SubagentType)
	}
	if n.CapabilityMode != CapReadOnly {
		t.Errorf("cap=%q", n.CapabilityMode)
	}
	if n.Isolation != IsoNone {
		t.Errorf("iso=%q", n.Isolation)
	}
	if n.Thoroughness != ThoroughnessMedium {
		t.Errorf("thorough=%q", n.Thoroughness)
	}
}

func TestNormalizeGeneralDefaultsAll(t *testing.T) {
	n, err := (SpawnRequest{Prompt: "implement X", SubagentType: "general"}).Normalize()
	if err != nil {
		t.Fatal(err)
	}
	if n.SubagentType != TypeGeneralPurpose || n.CapabilityMode != CapAll {
		t.Fatalf("got type=%q cap=%q", n.SubagentType, n.CapabilityMode)
	}
}

func TestNormalizeMutualExclusionCWDWorktree(t *testing.T) {
	err := (SpawnRequest{
		Prompt:    "x",
		Isolation: "worktree",
		CWD:       "/tmp/other",
	}).Validate()
	if err == nil || !strings.Contains(err.Error(), "mutually exclusive") {
		t.Fatalf("want mutual exclusion error, got %v", err)
	}
}

func TestNormalizeResumeWithoutPrompt(t *testing.T) {
	n, err := (SpawnRequest{ResumeFrom: "sub-abc"}).Normalize()
	if err != nil {
		t.Fatal(err)
	}
	if n.ResumeFrom != "sub-abc" {
		t.Fatalf("resume=%q", n.ResumeFrom)
	}
}

func TestNormalizeRequiresPromptOrResume(t *testing.T) {
	if err := (SpawnRequest{}).Validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestNormalizeThoroughnessOnlyExplore(t *testing.T) {
	err := (SpawnRequest{
		Prompt:       "x",
		SubagentType: "plan",
		Thoroughness: "quick",
	}).Validate()
	if err == nil || !strings.Contains(err.Error(), "thoroughness") {
		t.Fatalf("want thoroughness error, got %v", err)
	}
}

func TestNormalizeExploreThoroughness(t *testing.T) {
	n, err := (SpawnRequest{
		Prompt:       "x",
		SubagentType: "explore",
		Thoroughness: "very-thorough",
	}).Normalize()
	if err != nil {
		t.Fatal(err)
	}
	if n.Thoroughness != ThoroughnessVeryThorough {
		t.Fatalf("got %q", n.Thoroughness)
	}
}
