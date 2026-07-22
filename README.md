# hawk-core-contracts

Shared contracts for the hawk ecosystem.

This repo holds stable cross-repo type definitions used by every engine and
`hawk` itself: severity levels, findings, tool contracts, event models, policy
verdicts, review results, verification reports, and agent session state.

**Tagline:** Shared contracts for the hawk ecosystem.

## Install

```sh
go get github.com/GrayCodeAI/hawk-core-contracts
```

## Quick Reference

| Package | Key Types | Purpose |
|---------|-----------|---------|
| `types/` | `Severity`, `Finding`, `FindingSlice`, `ParseSeverity` | Severity levels, findings, shared result vocabulary |
| `tools/` | `ToolCall`, `ToolResult`, `ToolMeta`, `FinalizeResult`, version lifecycle types | Provider-neutral tool call, result, identity, and finalization contracts |
| `events/` | `ToolEvent`, `TraceEvent`, `UsageInfo` | Normalized tool and trace event contracts |
| `policy/` | `Risk`, `PermissionVerdict`, `Allow`, `Deny` | Risk, permission verdict, guardian decision, approval request contracts |
| `review/` | `Result`, `Finding`, `Comment` | Neutral review findings, comments, stats, and result contracts |
| `verify/` | `Report`, `Finding` | Neutral verification findings, stats, and report contracts |
| `sessions/` | `Phase`, `CostAccumulator`, `ParsePhase` | Cross-repo agent session state types |
| `agent/` | `SpawnRequest`, `SpawnResult`, hook event names | Typed subagent spawn + hook vocabulary |
| `llm/` | `Provider` + roles, `EventStreamer`, `EyrieMessage`, `EyrieResponse`, `ChatOptions`, `EyrieStreamEvent`, `StreamResult`, `Model`, `Usage`, … | Canonical provider port contract — the hawk↔eyrie boundary (`ToolCall`/`ToolResult` alias `tools/`) |

## Architecture

```
hawk-core-contracts (stdlib only)
├── types/     Severity, Finding, FindingSlice — the core vocabulary
├── tools/     ToolCall, ToolResult — provider-neutral tool contracts
├── events/    ToolEvent, TraceEvent — normalized event contracts
├── policy/    Risk, PermissionVerdict — governance contracts
├── review/    Result, Finding, Comment — review result contracts
├── verify/    Report, Finding — verification report contracts
├── sessions/  Phase, CostAccumulator — session state contracts
├── agent/     SpawnRequest, SpawnResult, hook events — subagent spawn contracts
└── llm/       Provider + 7 roles, EventStreamer, conversation DTOs — hawk↔eyrie port
```

`llm.ToolCall` / `llm.ToolResult` are type aliases of `tools.ToolCall` /
`tools.ToolResult` (single vocabulary). Host streaming on the port is
pull-based (`EventStreamer`); channel-based `StreamResult` remains for
lower-level client transports.


## Scope

### Allowed here

- shared enums (severity levels, risk levels, phases)
- shared structs (findings, tool calls, events, verdicts, reports)
- event models (tool events, trace events, usage info)
- finding/result models (severity, confidence, status)
- engine request/response contracts
- policy and tool contracts

### Not allowed here

- CLI code (commands, flags, shell output)
- provider implementations
- runtime logic
- storage implementations
- product orchestration
- anything that belongs in a single consuming repo

Engines should depend on this repo only when they produce or consume a shared
cross-repo contract. Eyrie depends on `llm/` for the host port; other engines
(e.g. `yaad`, `trace`) stay contract-free unless they share a real cross-repo
type.

If a type is only used inside one repo, it should stay in that repo.

## Migration History

The legacy `github.com/GrayCodeAI/hawk/shared/types` package has been removed.
All shared finding and severity definitions now live here. Migration is
complete:

1. Severity and finding definitions migrated from `hawk/shared/types`
2. `sight` and `inspect` migrated to import this repo
3. Tool, event, and policy contracts added
4. Review and verification result contracts added
5. Sessions package added for cross-repo session state
6. Tool contract versioning & finalization: `ToolMeta` identity envelope, additive-only
   version bump rule, closed `ToolNamespace` enum, `BehaviorPreset` back-compat presets,
   and the single-locking `FinalizeToolConfig` lifecycle (with `FinalizeResult`,
   `VersionWarning`, `FinalizeConfigViolation`)

## Ecosystem

hawk-core-contracts is a **foundation repo** in the hawk-eco mono-ecosystem:

| Component | Purpose |
|-----------|---------|
| **hawk-core-contracts** | Shared cross-repo contracts (this repo) |
| **hawk-mcpkit** | Shared MCP server scaffolding |
| **eyrie** | LLM provider runtime — routing, streaming, retries, caching |
| **yaad** | Graph-based persistent memory for coding agents |
| **tok** | Tokenizer, compression, secrets scanning, rate limiting |
| **sight** | Diff-based code review and static analysis |
| **inspect** | Security audit library (CVE, API security, CI output) |
| **trace** | Session capture and replay CLI |
| **hawk** | AI coding agent (this repo) |

`hawk` and all engines import `hawk-core-contracts` when they share a real
cross-repo contract; the repo itself never imports back.

## Tool Contract Versioning & Finalization

Tool configuration has a single locking **finalize** step instead of a scattered set of
enable/disable/options calls that can leave the tool set in an inconsistent state.

**`ToolMeta`** is the canonical identity envelope attached to every tool-call event. Its
`version` string is **additive-only**: new additive fields do not change it; only breaking
changes (field removal, type change, reserved-range change) bump it. `namespace` is a
**closed enum** (`ToolNamespace`) — a new unknown namespace intentionally fails strict typed
deserialization (`ToolNamespaceFrom` returns an error), so a forward-rolled contract can't
silently mis-route tools to a harness that doesn't understand them. `label` is a
cross-harness grouping key; `read_only` marks tools the harness must treat as side-effect free.

**`FinalizeToolConfig`** is the single, locking RPC that atomically commits a
`behavior_version` (a `BehaviorPreset`: `current`, `legacy`, or unspecified) plus an
enabled-tool set, and returns any `VersionWarning`s (non-fatal drift) or
`FinalizeConfigViolation`s (fatal, deterministic failures). After a call where
`FinalizeResult.Ok()` returns true (`Finalized && no violations`), tools are safe to call.

This models the lifecycle in SpaceXAI `grok`, but starts clean — the old scattered
Enable/Disable/Get/Set-ToolOptions surface is intentionally absent.

## Ecosystem Boundaries

Rules that keep this repo at the foundation layer:

- **Zero hawk-eco dependencies.** This repo imports only the Go standard
  library. `make boundaries` (also run in CI) enforces this with
  `scripts/check-ecosystem-boundaries.sh`.
- **Implementation-free.** See Scope above — no CLI code, provider
  implementations, runtime logic, storage, or orchestration.
- **Consumers, not dependents.** `hawk` and engines import this repo when
  they share a real cross-repo contract; it never imports them back.

If a change here would require importing anything outside the standard
library, that type does not belong in this repo.

## Codegen

`proto/hawk/contracts/v1/*.proto` mirrors every exported type above, one
`.proto` file per Go package. It exists for two things the hand-written Go
package alone can't give you:

- **Schema-level breaking-change detection.** CI's `proto` job runs
  `buf breaking` against `main` on every PR — catches the class of bug
  already seen once in this ecosystem (`hawk-sdk-go` independently
  hand-rolled its own `ToolResult` with a field named `tool_call_id`,
  diverging from this repo's `tool_use_id`, because nothing checked for it).
- **Python / TypeScript codegen**, for whenever a non-Go consumer wants this
  vocabulary instead of hand-porting it. `buf generate` produces
  `gen/go/`, `gen/python/`, `gen/typescript/` — not committed (regenerate
  with `make proto`), and not currently imported by anything: `hawk-sdk-go`,
  `hawk-sdk-python`, and `yaad`'s TypeScript SDK all still hand-port their
  own subset of this vocabulary today (`hawk-sdk-python/src/hawk/sessions.py`
  and `types.py` are the main example) and haven't adopted the generated
  packages. That's a real adoption decision for those repos to make on
  their own schedule, not something this repo forces.

**The `.proto` files and the Go structs are two independent, hand-kept-in-sync
definitions** — `gen/go/` is deliberately its own nested Go module (see
`gen/go/go.mod`) precisely so depending on `google.golang.org/protobuf`
there never touches this repo's zero-dependency root module. When you add
or change an exported Go type, update the matching `.proto` message in the
same PR — see `AGENTS.md`.

Some fields don't map 1:1 by protobuf convention; each divergence is
commented in the `.proto` source at the point it occurs (e.g. `Severity`'s
zero value is `SEVERITY_INFO`, not an `_UNSPECIFIED` sentinel, to keep its
numeric values identical to the Go `iota` constants that
`map[Severity]int` fields serialize by).

## Package Ownership

| Path | Team |
|------|------|
| `/types/` | `@GrayCodeAI/llm-team` |
| `/tools/` | `@GrayCodeAI/llm-team` |
| `/events/` | `@GrayCodeAI/llm-team` |
| `/policy/` | `@GrayCodeAI/llm-team` |
| `/review/` | `@GrayCodeAI/llm-team` |
| `/verify/` | `@GrayCodeAI/llm-team` |
| `/sessions/` | `@GrayCodeAI/llm-team` |
| `/agent/` | `@GrayCodeAI/llm-team` |
| `/llm/` | `@GrayCodeAI/llm-team` |
| `/proto/` | `@GrayCodeAI/llm-team` |
| `/VERSION` | `@GrayCodeAI/maintainers` |
| `/Makefile` | `@GrayCodeAI/devops-team` |
| `/*.md` | `@GrayCodeAI/docs-team` |
| `/.github/` | `@GrayCodeAI/devops-team` |

## License

MIT — see [LICENSE](LICENSE).
