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
| `tools/` | `ToolCall`, `ToolResult` | Provider-neutral tool call and result contracts |
| `events/` | `ToolEvent`, `TraceEvent`, `UsageInfo` | Normalized tool and trace event contracts |
| `policy/` | `Risk`, `PermissionVerdict`, `Allow`, `Deny` | Risk, permission verdict, guardian decision, approval request contracts |
| `review/` | `Result`, `Finding`, `Comment` | Neutral review findings, comments, stats, and result contracts |
| `verify/` | `Report`, `Finding` | Neutral verification findings, stats, and report contracts |
| `sessions/` | `Phase`, `CostAccumulator`, `ParsePhase` | Cross-repo agent session state types |

## Architecture

```
hawk-core-contracts (stdlib only)
├── types/     Severity, Finding, FindingSlice — the core vocabulary
├── tools/     ToolCall, ToolResult — provider-neutral tool contracts
├── events/    ToolEvent, TraceEvent — normalized event contracts
├── policy/    Risk, PermissionVerdict — governance contracts
├── review/    Result, Finding, Comment — review result contracts
├── verify/    Report, Finding — verification report contracts
└── sessions/  Phase, CostAccumulator — session state contracts
```

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
cross-repo contract. Contract-free engines (e.g., `eyrie`, `yaad`, `trace`)
should not add the dependency just for consistency.

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
| `/VERSION` | `@GrayCodeAI/maintainers` |
| `/Makefile` | `@GrayCodeAI/devops-team` |
| `/*.md` | `@GrayCodeAI/docs-team` |
| `/.github/` | `@GrayCodeAI/devops-team` |

## License

MIT — see [LICENSE](LICENSE).
