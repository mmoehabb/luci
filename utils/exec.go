package utils

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mmoehabb/luci/types"
)

func Act(c types.Config, inputs []string) {
	action := Dig(c, inputs)

	if action == nil {
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

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Panicf("cmd.CombinedOutput() failed with %s\n", err)
		}
		fmt.Printf("\n%s\n", string(output))
		return true

	case []string:
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.Command("cmd", "/C", strings.Join(action, " && "))
		} else {
			cmd = exec.Command("/bin/sh", "-c", strings.Join(action, " && "))
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Panicf("cmd.CombinedOutput() failed with %s\n", err)
		}
		fmt.Printf("\n%s\n", string(output))
		return true
	}

	return false
}
