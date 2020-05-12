package cli

import "fmt"

var (
	appVersion = "0"
	gitCommit  = "n/a"
	goVersion  = "n/a"
)

func versionString() string {
	return fmt.Sprintf(
		"coda-indexer v%s (commit: %s, go: %s)",
		appVersion,
		gitCommit,
		goVersion,
	)
}
