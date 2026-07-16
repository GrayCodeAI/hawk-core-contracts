package agent

import "testing"

func TestCanonicalHookEvent(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{HookPreToolUse, HookPreToolUse},
		{"pre_tool_use", HookPreToolUse},
		{"PreToolUse", HookPreToolUse},
		{"subagent_start", HookSubagentStart},
		{"on_error", HookFailure},
		{"", ""},
		{"not-a-hook", ""},
	}
	for _, tc := range cases {
		if got := CanonicalHookEvent(tc.in); got != tc.want {
			t.Errorf("CanonicalHookEvent(%q)=%q want %q", tc.in, got, tc.want)
		}
	}
}
