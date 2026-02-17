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
func Dig(action any, inputs []string) any {
	for i, input := range inputs {
		switch actTyped := action.(type) {
		case types.ShellConfig:
			action = actTyped[input]

		case types.AnnotatedAction:
			action = Dig(actTyped.Value, inputs[i:])
			continue

		case map[string]any:
			if actTyped["value"] != nil {
				action = Dig(MapToAnnotatedAction(actTyped), inputs[i:])
				continue
			}
			action = actTyped[input]
		}
	}
	return action
}
