// Package version provides version information for the application.
// This package can be imported by external tools if needed.
package version

import (
	"fmt"
	"runtime"
)

// Build information. These variables are populated at build time using -ldflags.
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = ""
)

// BuildInfo represents the build information
type BuildInfo struct {
	Version   string `json:"version"`
	BuildTime string `json:"buildTime"`
	GitCommit string `json:"gitCommit,omitempty"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

// GetBuildInfo returns the build information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		BuildTime: BuildTime,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func String() string {
	info := GetBuildInfo()
	if info.GitCommit != "" {
		return fmt.Sprintf("%s (commit: %s, built: %s, go: %s, platform: %s)",
			info.Version, info.GitCommit, info.BuildTime, info.GoVersion, info.Platform)
	}
	return fmt.Sprintf("%s (built: %s, go: %s, platform: %s)",
		info.Version, info.BuildTime, info.GoVersion, info.Platform)
}
