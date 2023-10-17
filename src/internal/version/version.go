package version

import (
	_ "embed"
)

var (
	// Version is the version of the turnkey CLI.
	//go:embed data/version
	Version string

	// Commit is the commit SHA on which this version of the CLI tool was built.
	//go:embed data/commit
	Commit  string
)
