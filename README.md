# hawk-core-contracts

Shared contracts for the Hawk ecosystem.

This repo exists to hold stable cross-repo definitions used by:

- `hawk`
- `eyrie`
- `yaad`
- `tok`
- `trace`
- `sight`
- `inspect`
- Hawk SDKs and extension surfaces where needed

## Scope

Allowed here:

- shared enums
- shared structs
- event models
- finding/result models
- engine request/response contracts
- policy and tool contracts

Not allowed here:

- CLI code
- provider implementations
- runtime logic
- storage implementations
- product orchestration

## Migration history

The legacy `github.com/GrayCodeAI/hawk/shared/types` package has been removed.
Severity, findings, and the packages below are the supported cross-repo API.

Engines should depend on this repo only when they produce or consume a shared
contract. Contract-free engines (for example `eyrie`, `yaad`, `trace`) should
not add the dependency just for consistency.

## Package map

- `types/` - severity, findings, and shared result vocabulary
- `tools/` - provider-neutral tool call and tool result contracts
- `events/` - normalized tool and trace event contracts
- `policy/` - risk, permission verdict, guardian decision, approval request contracts
- `review/` - neutral review findings, comments, stats, and result contracts
- `verify/` - neutral verification findings, stats, and report contracts

## Current status

Completed:

1. shared finding and severity definitions moved here
2. `sight` and `inspect` migrated to import this repo
3. Hawk docs and READMEs updated
4. tool, event, and policy contracts added
5. review and verification result contracts added

## Governance rules

- keep this repo implementation-free
- prefer additive changes
- avoid product-specific runtime assumptions
- do not move Hawk orchestration code here
- if a type is only used inside one repo, it should stay in that repo

## Ecosystem Boundaries

`hawk-core-contracts` is a **foundation repo** in the [hawk ecosystem](https://github.com/GrayCodeAI/hawk/blob/main/docs/architecture/hawk-ecosystem-summary.md) —
it sits below every engine and below `hawk` itself, alongside
`hawk-mcpkit`.

Rules that keep it there:

- **Zero hawk-eco dependencies.** This repo must never import `hawk`, any
  engine (`eyrie`, `yaad`, `tok`, `trace`, `sight`, `inspect`), any SDK, or
  `hawk-mcpkit`. Only the Go standard library. `make boundaries` (also run
  in CI) enforces this with `scripts/check-ecosystem-boundaries.sh`.
- **Implementation-free.** See Scope above — no CLI code, provider
  implementations, runtime logic, storage, or orchestration.
- **Consumers, not dependents.** `hawk` and engines import this repo when
  they share a real cross-repo contract; it never imports them back.

If a change here would require importing anything outside the standard
library, that type does not belong in this repo.
