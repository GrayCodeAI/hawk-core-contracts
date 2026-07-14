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
```

## Scope Rules

**Allowed:** shared enums, structs, event models, policy/tool contracts
**Not allowed:** CLI code, provider implementations, runtime logic, storage, product orchestration

## Ecosystem Boundaries

- Zero hawk-eco dependencies — stdlib only
- Implementation-free: no CLI, no providers, no runtime, no storage
- Consumers import this repo; it never imports them back

For full hawk-eco extension guidelines, see [hawk/AGENTS.md](https://github.com/GrayCodeAI/hawk/blob/main/AGENTS.md).
