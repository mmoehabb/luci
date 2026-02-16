package utils

import (
	"runtime"

	"github.com/mmoehabb/luci/types"
)

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

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
