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
