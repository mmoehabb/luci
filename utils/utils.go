package utils

import (
	"runtime"

	"github.com/mmoehabb/luci/types"
)

// GetShellType detects and returns the appropriate shell type for the current
// operating system. It returns Bash for Linux, Zshell for macOS (darwin),
// Powershell for Windows, and Unknown for any other OS.
func GetShellType() types.ShellType {
	switch runtime.GOOS {
	case "linux":
		return types.Bash
	case "darwin":
		return types.Zshell
	case "windows":
		return types.Powershell
	default:
		return types.Unknown
	}
}

// GetShellConfig returns a pointer to the shell-specific configuration section
// from the provided config based on the current operating system's shell type.
// It selects between Bash, Zshell, or Powershell configurations, defaulting to Bash
// if the shell type is unknown.
func GetShellConfig(c types.Config) *types.ShellConfig {
	var shellConfig *types.ShellConfig

	switch GetShellType() {
	case types.Bash:
		shellConfig = &c.Bash
	case types.Zshell:
		shellConfig = &c.Zshell
	case types.Powershell:
		shellConfig = &c.Powershell
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

// GetWildcardConfig returns a pointer to the wildcard ("*") configuration section.
// This section acts as a fallback when an action is not found in the shell-specific config.
func GetWildcardConfig(c types.Config) *types.ShellConfig {
	return &c.Wildcard
}

// GetMergedShellConfig returns a pointer to a new ShellConfig that deep-merges
// the shell-specific config with the wildcard ("*") config. The shell config
// takes priority on key conflicts at every nesting level.
func GetMergedShellConfig(c types.Config) *types.ShellConfig {
	shell := *GetShellConfig(c)
	wildcard := *GetWildcardConfig(c)
	merged := MergeShellConfigs(shell, wildcard)
	return &merged
}

// MergeShellConfigs deep-merges overlay into base. Base takes priority on key
// conflicts. When both sides have a map[string]any at the same key, the merge
// recurses; otherwise base wins.
func MergeShellConfigs(base, overlay types.ShellConfig) types.ShellConfig {
	result := make(types.ShellConfig)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range overlay {
		if existing, exists := result[k]; exists {
			existingMap, eOk := existing.(map[string]any)
			overlayMap, oOk := v.(map[string]any)
			if eOk && oOk {
				result[k] = mergeMaps(existingMap, overlayMap)
			}
		} else {
			result[k] = v
		}
	}
	return result
}

func mergeMaps(base, overlay map[string]any) map[string]any {
	result := make(map[string]any)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range overlay {
		if existing, exists := result[k]; exists {
			existingMap, eOk := existing.(map[string]any)
			overlayMap, oOk := v.(map[string]any)
			if eOk && oOk {
				result[k] = mergeMaps(existingMap, overlayMap)
			}
		} else {
			result[k] = v
		}
	}
	return result
}
