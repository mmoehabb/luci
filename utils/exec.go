package utils

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/mmoehabb/luci/types"
)

// Act is the main entry point for executing actions based on user input.
// It takes a configuration and a list of input arguments, resolves the appropriate
// action from the shell configuration (with wildcard fallback), and executes it.
// If no matching action is found, it prints usage information or displays available actions.
func Act(c types.Config, inputs []string) {
	merged := *GetMergedShellConfig(c)
	action, i := Dig(merged, inputs)
	args := inputs[i:]

	if action == nil {
		if len(inputs) <= 1 {
			PrintUsage(c)
			return
		}
		PrintActionWithInputs(merged, inputs[0:len(inputs)-1], 0)
		return
	}

	switch action.(type) {
	case types.AnnotatedAction:
		val := action.(types.AnnotatedAction).Value

		switch val.(type) {
		case string:
			executed := execAction(val, args)
			if executed == true {
				return
			}
		case []string:
			executed := execAction(val, args)
			if executed == true {
				return
			}
		}

	case []string:
		executed := execAction(action, args)
		if executed == true {
			return
		}

	case string:
		executed := execAction(action, args)
		if executed == true {
			return
		}
	}

	PrintActionWithInputs(merged, inputs, 0)
}

func execAction(action any, args []string) bool {
	switch action := action.(type) {
	case string:
		var cmd *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/C", action+" "+strings.Join(args, " "))
		case "darwin":
			cmd = exec.Command("/bin/zsh", "-c", action+" "+strings.Join(args, " "))
		default:
			cmd = exec.Command("/bin/sh", "-c", action+" "+strings.Join(args, " "))
		}

		PrintCommand(action, args)
		execCmd(cmd)
		return true

	case []string:
		var cmd *exec.Cmd
		var cmdStr = strings.Join(action, " && ")

		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd", "/C", cmdStr+" "+strings.Join(args, " "))
		case "darwin":
			cmd = exec.Command("/bin/zsh", "-c", cmdStr+" "+strings.Join(args, " "))
		default:
			cmd = exec.Command("/bin/sh", "-c", cmdStr+" "+strings.Join(args, " "))
		}

		PrintCommand(cmdStr, args)
		execCmd(cmd)
		return true
	}

	return false
}

func execCmd(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		color.Red("Failed to start command: %s\n", err)
		os.Exit(1)
	}

	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		color.Red("Command failed: %s\n", err)
		os.Exit(1)
	}
}
