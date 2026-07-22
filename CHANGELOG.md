# Changelog

All notable changes to `hawk-core-contracts` are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed

- `VERSION` file now reports `0.1.7` (was left at `0.1.6` after the llm
  package release; `version.go` embeds this file).
- `llm.ToolCall` / `llm.ToolResult` are type aliases of `tools.ToolCall` /
  `tools.ToolResult` so the ecosystem has one tool-call vocabulary.
- Host-facing `llm` types aligned with the eyrie engine facade:
  `CatalogHealth`, `DeploymentSummary`, `PreflightReport` (`LiveVerified`),
  `ProviderStateSecurity`, and pull-based `EventStreamer` on `Generator.Stream`.
- `proto/hawk/contracts/v1/llm.proto` mirrors llm conversation / host DTOs;
  `make proto` regenerates agent + llm stubs under `gen/`.
- Unit tests for `llm` (`StreamResult`, aliases, constants, `EventStreamer`).

## [0.1.7] — 2026-07-22

### Added

- `llm/` package: shared LLM provider contract types (`provider.go`,
  `types.go`) ported from the engine surface, consumable by downstream
  repos (e.g. `eyrie`) without a local replace directive.

## [0.1.6] — 2026-07-16

### Added

- `agent/` package: typed subagent spawn contracts (`SpawnRequest`,
  `SpawnResult`, capability/isolation/subagent-type parse helpers,
  `Normalize`/`Validate`) and canonical hook event names with vendor
  aliases (`CanonicalHookEvent`). Year 0 Grok behavioral port (PACK-01).
- `proto/hawk/contracts/v1/agent.proto` mirroring `SpawnRequest` /
  `SpawnResult` for buf lint/breaking and cross-language codegen.

## [0.1.5] — 2026-07-11

### Added

- `proto/` — a Buf-managed protobuf schema mirroring every exported Go type,
  used for schema-level breaking-change detection (`buf breaking`, gated in
  CI on every PR once this lands) and for generating Python/TypeScript
  clients (`buf generate` → `gen/python/`, `gen/typescript/`; not committed,
  regenerate on demand). The hand-written Go package is unaffected and stays
  the source of truth for Go consumers — see `proto/README.md` (or the repo
  README's Codegen section) for the mapping rules between the two.
- CI: `proto` job runs `buf lint` and `buf breaking` (against `main`) on
  every PR that touches `proto/`.

No contract package changes since `v0.1.4` — this release line so far is
CI/tooling and the addition above.

## [0.1.4] — 2026-07-11

### Added

- `VERSION` file + `//go:embed`-based `contracts.Version`, replacing ad hoc
  versioning with a single source of truth (#6).
- CI: race-detector + coverage-threshold gate.

### Changed

- `policy.ParseRisk` now uses `strings.ToLower`/`strings.TrimSpace` instead
  of a hand-rolled lowercasing helper — internal only, no behavior change.
- Go version bumped to 1.26.5.

No wire-format / contract-shape changes in this release.

## [0.1.3] — 2026-07-05

Added the MIT `LICENSE` file. No contract changes.

## [0.1.2] — 2026-07-04

### Added

- **New `sessions/` package**: `Phase`, `SessionID`, `ContextSnapshot`,
  `ToolCallRecord`, `PhaseUsage`, `CostAccumulator` — cross-repo agent
  session identity, pipeline-phase tagging, and cost accounting, so engines
  and hawk's orchestrator share one vocabulary for this without a circular
  dependency.
- `CODEOWNERS`, CI governance workflow.

## [0.1.1] — 2026-06-25

### Added

- `review.SASTFusionResult` (`Confirmed`/`Dismissed`/`Unaddressed` finding
  lists) on `review.Result`, populated when SAST-LLM fusion is active (#1).

## [0.1.0] — 2026-06-21

Initial governed release. Establishes repo governance (`Makefile`, boundary
guard, engine-grade CI) matching the rest of the hawk-eco foundation/engine
repos.

Contract packages at this version — six packages; `sessions/` did not exist
yet (added in `0.1.2`, see above):

- `types/` — severity, findings, shared result vocabulary
- `tools/` — provider-neutral tool call and tool result contracts
- `events/` — normalized tool and trace event contracts
- `policy/` — risk, permission verdict, guardian decision, approval request contracts
- `review/` — neutral review findings, comments, stats, and result contracts
- `verify/` — neutral verification findings, stats, and report contracts
