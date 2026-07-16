package tools_test

import (
	"testing"

	"github.com/GrayCodeAI/hawk-core-contracts/tools"
)

func TestBehaviorPresetParse(t *testing.T) {
	valid := map[string]tools.BehaviorPreset{
		"current": tools.BehaviorPresetCurrent,
		"legacy":  tools.BehaviorPresetLegacy,
		"":        tools.BehaviorPresetUnspecified,
	}
	for in, want := range valid {
		got, err := tools.BehaviorPresetFrom(in)
		if err != nil {
			t.Fatalf("BehaviorPresetFrom(%q) error = %v", in, err)
		}
		if got != want {
			t.Fatalf("BehaviorPresetFrom(%q) = %q, want %q", in, got, want)
		}
	}

	invalid := []string{"future", "CURRENT", "Current", "legacy-v2"}
	for _, in := range invalid {
		got, err := tools.BehaviorPresetFrom(in)
		if err == nil {
			t.Fatalf("BehaviorPresetFrom(%q) = %q, want error", in, got)
		}
		if got != tools.BehaviorPresetUnspecified {
			t.Fatalf("BehaviorPresetFrom(%q) zero value = %q, want unspecified", in, got)
		}
	}
}

func TestFinalizeResultOk(t *testing.T) {
	cases := []struct {
		name string
		in   tools.FinalizeResult
		want bool
	}{
		{"not finalized", tools.FinalizeResult{Finalized: false}, false},
		{"finalized with violation", tools.FinalizeResult{
			Finalized:  true,
			Violations: []tools.FinalizeConfigViolation{{Code: tools.FinalizeErrorCodeUnknownTool, ToolID: "x"}},
		}, false},
		{"finalized with warning only", tools.FinalizeResult{
			Finalized: true,
			Warnings:  []tools.VersionWarning{{Code: tools.FinalizeErrorCodeDeprecatedReplaced}},
		}, true},
		{"finalized clean", tools.FinalizeResult{Finalized: true, BehaviorVersion: "current"}, true},
	}
	for _, tc := range cases {
		if got := tc.in.Ok(); got != tc.want {
			t.Fatalf("%s: Ok() = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestToolNamespaceFrom_ClosedEnum(t *testing.T) {
	known := []string{"", "hawk_build", "hawk_build_concise", "codex", "opencode", "mcp", "acp"}
	for _, in := range known {
		got, err := tools.ToolNamespaceFrom(in)
		if err != nil {
			t.Fatalf("ToolNamespaceFrom(%q) error = %v", in, err)
		}
		if string(got) != in {
			t.Fatalf("ToolNamespaceFrom(%q) = %q", in, got)
		}
	}

	// Forward-safety: an unknown namespace must error, not silently map.
	for _, in := range []string{"future_ns", "grok_build", "HAWK_BUILD", "Mcp"} {
		got, err := tools.ToolNamespaceFrom(in)
		if err == nil {
			t.Fatalf("ToolNamespaceFrom(%q) = %q, want error (closed enum)", in, got)
		}
		if got != tools.ToolNamespaceUnspecified {
			t.Fatalf("ToolNamespaceFrom(%q) zero value = %q, want unspecified", in, got)
		}
	}
}

func TestToolMetaFields(t *testing.T) {
	meta := tools.ToolMeta{
		Version:   "1",
		Name:      "exec",
		Kind:      "shell",
		Namespace: tools.ToolNamespaceHawkBuild,
		Label:     "build-tools",
		ReadOnly:  true,
	}

	if meta.Version != "1" {
		t.Fatalf("Version = %q", meta.Version)
	}
	if meta.Name != "exec" {
		t.Fatalf("Name = %q", meta.Name)
	}
	if meta.Kind != "shell" {
		t.Fatalf("Kind = %q", meta.Kind)
	}
	if meta.Namespace != tools.ToolNamespaceHawkBuild {
		t.Fatalf("Namespace = %q", meta.Namespace)
	}
	if meta.Label != "build-tools" {
		t.Fatalf("Label = %q", meta.Label)
	}
	if !meta.ReadOnly {
		t.Fatalf("ReadOnly = false, want true")
	}
}
