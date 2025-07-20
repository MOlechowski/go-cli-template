package version

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBuildInfo(t *testing.T) {
	// Save original values
	origVersion := Version
	origBuildTime := BuildTime
	origGitCommit := GitCommit
	
	// Set test values
	Version = "1.2.3"
	BuildTime = "2024-01-01T00:00:00Z"
	GitCommit = "abc123"
	
	defer func() {
		// Restore original values
		Version = origVersion
		BuildTime = origBuildTime
		GitCommit = origGitCommit
	}()
	
	info := GetBuildInfo()
	
	assert.Equal(t, "1.2.3", info.Version)
	assert.Equal(t, "2024-01-01T00:00:00Z", info.BuildTime)
	assert.Equal(t, "abc123", info.GitCommit)
	assert.Equal(t, runtime.Version(), info.GoVersion)
	assert.Contains(t, info.Platform, runtime.GOOS)
	assert.Contains(t, info.Platform, runtime.GOARCH)
}

func TestString(t *testing.T) {
	// Save original values
	origVersion := Version
	origBuildTime := BuildTime
	origGitCommit := GitCommit
	
	defer func() {
		// Restore original values
		Version = origVersion
		BuildTime = origBuildTime
		GitCommit = origGitCommit
	}()
	
	t.Run("with git commit", func(t *testing.T) {
		Version = "1.2.3"
		BuildTime = "2024-01-01"
		GitCommit = "abc123"
		
		result := String()
		assert.Contains(t, result, "1.2.3")
		assert.Contains(t, result, "commit: abc123")
		assert.Contains(t, result, "built: 2024-01-01")
		assert.Contains(t, result, "go: "+runtime.Version())
	})
	
	t.Run("without git commit", func(t *testing.T) {
		Version = "1.2.3"
		BuildTime = "2024-01-01"
		GitCommit = ""
		
		result := String()
		assert.Contains(t, result, "1.2.3")
		assert.NotContains(t, result, "commit:")
		assert.Contains(t, result, "built: 2024-01-01")
	})
}

func TestDefaultValues(t *testing.T) {
	// Test that defaults are set correctly
	assert.NotEmpty(t, Version)
	assert.NotEmpty(t, BuildTime)
}