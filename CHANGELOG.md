# Changelog

All notable changes to `hawk-core-contracts` are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] — 2026-07-05

Initial governed release. Establishes `VERSION` as the source of truth and
adds repo governance (`CODEOWNERS`, `Makefile`, boundary guard, engine-grade
CI) matching the rest of the hawk-eco foundation/engine repos.

Existing contract packages at this version:

- `types/` — severity, findings, shared result vocabulary
- `tools/` — provider-neutral tool call and tool result contracts
- `events/` — normalized tool and trace event contracts
- `policy/` — risk, permission verdict, guardian decision, approval request contracts
- `review/` — neutral review findings, comments, stats, and result contracts
- `verify/` — neutral verification findings, stats, and report contracts
- `sessions/` — cross-repo agent session state types
