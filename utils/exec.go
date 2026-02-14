package utils

import (
	"bufio"
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
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(pipe)

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Wait for the command to finish after reading all output
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
