package config

import "fmt"

var (
	AppName    = "mina-indexer"
	AppVersion = "0.3.0.beta1"
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
