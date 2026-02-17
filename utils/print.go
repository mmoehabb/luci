package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mmoehabb/luci/types"
)

func PrintHeader(c types.Config) {
	color.HiGreen(`
	  /\\_/\\  
	 ( ^   ^ ) 
	  >  ^  < 
	`)
	color.HiGreen(`*** %s ***`, c.Title)

	colored := color.New(color.FgHiWhite).Sprint("> " + c.Description)
	wrapped := text.WrapSoft(colored, 60)
	fmt.Println(wrapped)

	color.Yellow("\nUsage:\n\n")
}

func PrintUsage(c types.Config) {
	PrintHeader(c)
	shell := *GetShellConfig(c)
	for action := range shell {
		PrintActionWithInputs(shell, []string{action}, 0)
	}
}

func PrintActionWithInputs(c map[string]any, inputs []string, level int) error {
	action := Dig(c, inputs)
	if action == nil {
		return errors.New("Action couldn't be found!")
	}
	PrintAction(action, inputs, level)
	return nil
}

func PrintAction(action any, inputs []string, level int) {
	switch action := action.(type) {
	case types.AnnotatedAction:
		color.New(color.FgMagenta).Printf("%sluci %s\n", indent(level), strings.Join(inputs, " "))
		if action.Title != "" {
			color.Blue("%s** %s **", indent(level+1), action.Title)
		}
		if action.Description != "" {
			color.Black("%s> %s", indent(level+1), action.Description)
		}
		switch annVal := action.Value.(type) {
		case map[string]any:
			if annVal["value"] != nil {
				fmt.Println("", inputs)
				PrintAction(MapToAnnotatedAction(annVal), inputs, level+1)
				return
			}
		}
		PrintAction(action.Value, inputs, level+1)

	case map[string]any:
		if action["value"] != nil {
			PrintAction(MapToAnnotatedAction(action), inputs, level)
			return
		}
		for key := range action {
			PrintAction(action[key], append(inputs, key), level)
		}

	case []string:
		fmt.Printf("%sluci %s\n", indent(level), strings.Join(inputs, " "))

	case string:
		fmt.Printf("%sluci %s\n", indent(level), strings.Join(inputs, " "))
	}
}

func PrintCommand(cmd string) {
	color.New(color.BgGreen, color.FgHiWhite).Printf("+ %s", cmd)
	fmt.Println()
}

func indent(count int) string {
	var res strings.Builder
	for range count {
		res.WriteString("\t")
	}
	return res.String()
}
