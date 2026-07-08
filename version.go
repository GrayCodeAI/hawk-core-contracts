// Package contracts provides cross-repo type contracts for hawk-eco.
//
// The Version variable is sourced from the VERSION file at the repo root.
package contracts

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var versionFile string

// Version of the hawk-core-contracts library. Single source of truth: VERSION file.
var Version = strings.TrimSpace(versionFile)
