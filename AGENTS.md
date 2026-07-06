# AGENTS.md — hawk-core-contracts

This file describes the hawk-core-contracts project for AI agents working in
this codebase. The TUI `/memory` command references this file.

---

## Project Overview

hawk-core-contracts is the shared cross-repo contract types package for the
hawk ecosystem. It holds stable definitions that every engine and hawk itself
depend on: severity levels, findings, tool contracts, event models, policy
verdicts, review results, verification reports, and agent session state.

**Tagline:** Shared contracts for the hawk ecosystem.

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
cross-repo contract. The repo itself never imports back.

## Architecture

```
hawk-core-contracts/
├── types/                  # Severity, findings, shared result vocabulary
│   ├── finding.go          # Finding, FindingSlice, FilterBySeverity
│   ├── severity.go         # Severity levels, ParseSeverity
│   └── *_test.go           # Table-driven tests
├── tools/                  # Provider-neutral tool call and result contracts
│   └── *.go
├── events/                 # Normalized tool and trace event contracts
│   └── *.go
├── policy/                 # Risk, permission verdict, guardian decision contracts
│   └── *.go
├── review/                 # Neutral review findings, comments, stats, results
│   └── *.go
├── verify/                 # Neutral verification findings, stats, reports
│   └── *.go
├── sessions/               # Cross-repo agent session state types
│   └── *.go
├── scripts/
│   └── check-ecosystem-boundaries.sh   # CI guard: zero hawk-eco deps
├── .github/
│   └── workflows/
│       ├── ci.yml                  # CI: format, module hygiene, vet, lint, test, security
│       └── release.yml             # GitHub Release on v* tags
├── Makefile                # Local dev tasks
├── lefthook.yml            # Pre-commit hooks (boundary guard, co-author strip)
├── AGENTS.md               # This file
├── README.md               # Package map, scope, governance rules
├── CHANGELOG.md            # Keep a Changelog format
├── CODEOWNERS              # Code ownership (package-level + infra)
├── LICENSE                 # MIT
├── VERSION                 # Source of truth for versioning
└── go.mod / go.sum         # Module files (stdlib only)
```

## Key Design Decisions

- **Implementation-free:** This repo holds only type definitions and
  constructors. No CLI code, no provider implementations, no runtime logic,
  no storage, no orchestration. See Scope below for the full list.
- **Zero hawk-eco dependencies:** This repo depends only on the Go standard
  library. `make boundaries` (also run in CI) enforces this. Violations are
  caught by `scripts/check-ecosystem-boundaries.sh`.
- **Additive-only:** Changes to contracts should be backward-compatible.
  New fields get zero values; new packages are added, never removed.
- **Prefixed consumers, not dependents:** `hawk` and engines import this
  repo when they share a real cross-repo contract; it never imports them
  back.
- **Single source of truth:** `VERSION` file is the canonical version
  identifier. `go.mod` uses `go 1.26.4`.

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

If a type is only used inside one repo, it should stay in that repo.

## Development Guidelines

### Build & Test

```bash
make test          # Run unit tests (no race detector)
make test-race     # Run unit tests with race detector
make lint          # Run golangci-lint
make boundaries    # Enforce zero hawk-eco imports
make ci            # Full CI suite (tidy, fmt, vet, boundaries, lint, test-race, security)
```

### Go Conventions

- Standard Go project layout: `cmd/` for entry points, `internal/` for private
- Tests live alongside source files (`foo.go` → `foo_test.go`)
- Package-level tests use `package <name>_test` to avoid importing internal
  packages when testing only public API surface
- Use table-driven tests where practical
- Errors are values — wrap with `fmt.Errorf("context: %w", err)`
- No global mutable state; prefer dependency injection

### Commit Conventions

Use [Conventional Commits](https://www.conventionalcommits.org/):
```
feat: add review result contracts
fix: handle nil findings in FilterBySeverity
refactor: extract severity parsing into shared helper
```

## Package Map

| Package | Contents | Owners |
|---------|----------|--------|
| `types/` | Severity levels, findings, shared result vocabulary | `@GrayCodeAI/llm-team` |
| `tools/` | Provider-neutral tool call and tool result contracts | `@GrayCodeAI/llm-team` |
| `events/` | Normalized tool and trace event contracts | `@GrayCodeAI/llm-team` |
| `policy/` | Risk, permission verdict, guardian decision, approval request contracts | `@GrayCodeAI/llm-team` |
| `review/` | Neutral review findings, comments, stats, and result contracts | `@GrayCodeAI/llm-team` |
| `verify/` | Neutral verification findings, stats, and report contracts | `@GrayCodeAI/llm-team` |
| `sessions/` | Cross-repo agent session state types | `@GrayCodeAI/llm-team` |

## Package Guidelines

Each contract package follows these conventions:

- **Pure data types:** structs, enums, constructors, and methods on those types
- **No external dependencies:** only stdlib imports
- **No I/O:** no file reads, HTTP calls, or database access
- **JSON-serializable:** all public structs should serialize cleanly via
  `encoding/json` (no `map[struct]` keys, no unexported fields without tags)
- **Nil-safety:** pointer receivers for nil-safe methods; document expected
  behavior for nil receivers
- **Test file naming:** `*_test.go` alongside source files
- **Package-level tests:** use `package <name>_test` when testing only the
  public API surface of a package

## Adding a New Contract Package

1. Create a new directory `<name>/` at the repo root
2. Define types, enums, and constructors in `.go` files
3. Add tests in `*_test.go` files (using `package <name>_test`)
4. Do not add any imports from the hawk ecosystem — only stdlib
5. Add the package to `CODEOWNERS` under the `@GrayCodeAI/llm-team` entry
6. Add any necessary types to the package map in the README
7. Add package-level docs ( exported type descriptions)

## Testing Patterns

- **Table-driven tests** with `t.Run(name, func(t *testing.T){...})` for all multi-case tests
- **`t.Parallel()`** on all tests that don't share mutable state
- **Package-level API tests** use `package <name>_test` to test only public surface without importing internal packages
- **No mocks framework** — use concrete types and test doubles
- **JSON round-trip tests** for serializable structs (marshal → unmarshal → compare fields)

## Common Pitfalls

- Do not import any hawk-eco package in this repo — only the Go standard library
- Do not put runtime logic, CLI code, or product orchestration here
- If a type is only used inside one repo, it should stay in that repo
- Keep changes additive; avoid renaming or restructuring existing types
- Do not add non-stdlib dependencies — if you need something, put it in the
  consuming engine instead
- `go.mod` has no `require` blocks (stdlib only); never add module dependencies

## File Organization Notes

| File | Purpose |
|------|---------|
| `types/` | Severity levels, findings, result vocabulary |
| `tools/` | Tool call/result contracts |
| `events/` | Tool and trace event contracts |
| `policy/` | Risk and permission verdict contracts |
| `review/` | Review result contracts |
| `verify/` | Verification report contracts |
| `sessions/` | Agent session state contracts |
| `scripts/check-ecosystem-boundaries.sh` | CI guard against hawk-eco imports |
| `.github/workflows/ci.yml` | CI pipeline |
| `.github/workflows/release.yml` | GitHub Release on `v*` tags |
| `Makefile` | Local dev tasks |
| `lefthook.yml` | Pre-commit hooks |
| `AGENTS.md` | This file |
| `README.md` | Package map, scope, governance rules |
| `CHANGELOG.md` | Keep a Changelog format |
| `CODEOWNERS` | Code ownership |
| `VERSION` | Canonical version |
