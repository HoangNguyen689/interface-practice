package version

import "fmt"

const unknown = "unknown"

var (
	gitCommit = unknown
	buildDate = unknown
	version   = unknown
)

type BuildInfo struct {
	Version   string
	GitCommit string
	BuildDate string
}

func (i BuildInfo) String() string {
	return fmt.Sprintf(
		"version: %s, git-commit: %s, build-date: %s",
		i.Version,
		i.GitCommit,
		i.BuildDate,
	)
}

func Get() BuildInfo {
	return BuildInfo{
		Version:   version,
		GitCommit: gitCommit,
		BuildDate: buildDate,
	}
}
