package utils

import (
	"runtime"

	"github.com/mmoehabb/luci/types"
)

// GetShellType detects and returns the appropriate shell type for the current
// operating system. It returns Bash for Linux, Zshell for macOS (darwin),
// Bat for Windows, and Unknown for any other OS.
func GetShellType() types.ShellType {
	switch runtime.GOOS {
	case "linux":
		return types.Bash
	case "darwin":
		return types.Zshell
	case "windows":
		return types.Bat
	default:
		return types.Unknown
	}
}

// GetShellConfig returns a pointer to the shell-specific configuration section
// from the provided config based on the current operating system's shell type.
// It selects between Bash, Zshell, or Bat configurations, defaulting to Bash
// if the shell type is unknown.
func GetShellConfig(c types.Config) *types.ShellConfig {
	var shellConfig *types.ShellConfig

	switch GetShellType() {
	case types.Bash:
		shellConfig = &c.Bash
	case types.Zshell:
		shellConfig = &c.Zshell
	case types.Bat:
		shellConfig = &c.Bat
	default:
		shellConfig = &c.Bash
	}

	return shellConfig
}

// Must panics if the provided error is not nil, otherwise it returns without error.
// This function is used for error handling in scenarios where errors indicate
// critical failures that should halt program execution.
func Must(err error) {
	if err != nil {
		panic(err)
	}
}
