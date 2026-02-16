package utils

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mmoehabb/luci/types"
)

func Act(c types.Config, inputs []string) {
	action := Dig(c, inputs)

	if action == nil {
		if len(inputs) <= 1 {
			PrintUsage(c)
			return
		}
		PrintAction(c, inputs[0:len(inputs)-1], 0)
		return
	}

	switch action.(type) {
	case types.AnnotatedAction:
		val := action.(types.AnnotatedAction).Value

		switch val.(type) {
		case string:
			executed := execAction(val)
			if executed == true {
				return
			}
		case []string:
			executed := execAction(val)
			if executed == true {
				return
			}
		}

	case []string:
		executed := execAction(action)
		if executed == true {
			return
		}

	case string:
		executed := execAction(action)
		if executed == true {
			return
		}
	}

	PrintAction(c, inputs, 0)
}

func execAction(action any) bool {
	switch action := action.(type) {
	case string:
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", action)
		} else {
			cmd = exec.Command("/bin/sh", "-c", action)
		}
		PrintCommand(action)
		execCmd(cmd)
		return true

	case []string:
		var cmd *exec.Cmd
		var cmdStr = strings.Join(action, " && ")
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", cmdStr)
		} else {
			cmd = exec.Command("/bin/sh", "-c", cmdStr)
		}
		PrintCommand(cmdStr)
		execCmd(cmd)
		return true
	}

	return false
}

func execCmd(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	Must(cmd.Start())
	Must(cmd.Wait())
}
