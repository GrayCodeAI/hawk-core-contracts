---
description: hawk-core-contracts — shared cross-repo types and boundaries.
globs: "*.go"
alwaysApply: false
---

# hawk-core-contracts Conventions

Shared contracts for the hawk ecosystem. Stdlib-only imports.

## Development workflow

When starting any new work (feature, fix, refactor, chore), always create a feature branch from `main` first. Never commit directly to `main`. Use branch naming conventions like `feat/<description>`, `fix/<description>`, or `chore/<description>`. Open a PR, ensure CI is green, then merge.

## Build & Test

```bash
go test ./...                    # Run tests
go vet ./...                     # Static analysis
go mod tidy                      # Tidy modules
make boundaries                  # Enforce ecosystem boundary rules
make proto                       # Lint proto/, check breaking changes, regenerate gen/
```

## Proto / Cross-Language Schema

`proto/` mirrors every exported Go type for schema-level breaking-change
detection and Python/TypeScript codegen (`gen/`, not committed — regenerate
with `make proto` / `buf generate`). It is additive, not authoritative: the
hand-written Go package above is still the source of truth for Go
consumers, and `gen/go/` lives in its own nested module
(`gen/go/go.mod`) specifically so pulling in `google.golang.org/protobuf`
there never touches the root module's zero-dependency guarantee.

**When you change an exported Go type, update the matching `.proto`
message in the same PR.** They are not generated from each other and can
drift silently otherwise. CI's `proto` job (`buf lint` + `buf breaking`
against `main`) will fail the build on an accidental breaking change to the
schema — if the break is intentional, explain why in the PR description and
add a CHANGELOG entry under a new version, same as any other breaking
change to this package.

## Scope Rules

**Allowed:** shared enums, structs, event models, policy/tool contracts
**Not allowed:** CLI code, provider implementations, runtime logic, storage, product orchestration

## Ecosystem Boundaries

- Zero hawk-eco dependencies — stdlib only
- Implementation-free: no CLI, no providers, no runtime, no storage
- Consumers import this repo; it never imports them back

For full hawk-eco extension guidelines, see [hawk/AGENTS.md](https://github.com/GrayCodeAI/hawk/blob/main/AGENTS.md).

<!-- gitnexus:start -->
## GitNexus — Code Intelligence

This project is indexed by GitNexus as **hawk-core-contracts** (476 symbols, 758 relationships, 4 execution flows). Use the GitNexus MCP tools to understand code, assess impact, and navigate safely.

> Index stale? Run `node .gitnexus/run.cjs analyze` from the project root — it auto-selects an available runner. No `.gitnexus/run.cjs` yet? `npx gitnexus analyze` (npm 11 crash → `npm i -g gitnexus`; #1939).

## Always Do

- **MUST run impact analysis before editing any symbol.** Before modifying a function, class, or method, run `impact({target: "symbolName", direction: "upstream"})` and report the blast radius (direct callers, affected processes, risk level) to the user.
- **MUST run `detect_changes()` before committing** to verify your changes only affect expected symbols and execution flows. For regression review, compare against the default branch: `detect_changes({scope: "compare", base_ref: "main"})`.
- **MUST warn the user** if impact analysis returns HIGH or CRITICAL risk before proceeding with edits.
- When exploring unfamiliar code, use `query({search_query: "concept"})` to find execution flows instead of grepping. It returns process-grouped results ranked by relevance.
- When you need full context on a specific symbol — callers, callees, which execution flows it participates in — use `context({name: "symbolName"})`.
- For security review, `explain({target: "fileOrSymbol"})` lists taint findings (source→sink flows; needs `analyze --pdg`).

## Never Do

- NEVER edit a function, class, or method without first running `impact` on it.
- NEVER ignore HIGH or CRITICAL risk warnings from impact analysis.
- NEVER rename symbols with find-and-replace — use `rename` which understands the call graph.
- NEVER commit changes without running `detect_changes()` to check affected scope.

## Resources

| Resource | Use for |
|----------|---------|
| `gitnexus://repo/hawk-core-contracts/context` | Codebase overview, check index freshness |
| `gitnexus://repo/hawk-core-contracts/clusters` | All functional areas |
| `gitnexus://repo/hawk-core-contracts/processes` | All execution flows |
| `gitnexus://repo/hawk-core-contracts/process/{name}` | Step-by-step execution trace |

## CLI

| Task | Read this skill file |
|------|---------------------|
| Understand architecture / "How does X work?" | `.claude/skills/gitnexus/gitnexus-exploring/SKILL.md` |
| Blast radius / "What breaks if I change X?" | `.claude/skills/gitnexus/gitnexus-impact-analysis/SKILL.md` |
| Trace bugs / "Why is X failing?" | `.claude/skills/gitnexus/gitnexus-debugging/SKILL.md` |
| Rename / extract / split / refactor | `.claude/skills/gitnexus/gitnexus-refactoring/SKILL.md` |
| Tools, resources, schema reference | `.claude/skills/gitnexus/gitnexus-guide/SKILL.md` |
| Index, status, clean, wiki CLI commands | `.claude/skills/gitnexus/gitnexus-cli/SKILL.md` |

<!-- gitnexus:end -->
This repo is a submodule of [hawk](https://github.com/GrayCodeAI/hawk) at `hawk/external/hawk-core-contracts`. Work in the submodule (go.work picks it up), push, sync here, PR/merge, then pull main in the submodule.
