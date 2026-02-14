package utils

import (
	"github.com/mmoehabb/luci/types"
)

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

// Return an Action or ActionRecord within the config c.
func Dig(c types.Config, inputs []string) any {
	shellc := *GetShellConfig(c)
	action := shellc[inputs[0]]
	if action == nil {
		return nil
	}

	for _, input := range inputs[1:] {
		switch action.(type) {
		case types.AnnotatedAction:
			val := action.(types.AnnotatedAction).Value
			switch val := val.(type) {
			case map[string]any:
				action = val[input]
			}

		case map[string]any:
			action = action.(map[string]any)[input]
		}
	}

	return action
}
