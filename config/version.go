package config

import "fmt"

var (
	AppName    = "mina-indexer"
	AppVersion = "0.9.0"
	GitCommit  = "-"
	GoVersion  = "-"
)

// VersionString returns the full app version string
func VersionString() string {
	return fmt.Sprintf(
		"%s %s (git: %s, %s)",
		AppName,
		AppVersion,
		GitCommit,
		GoVersion,
	)
}
