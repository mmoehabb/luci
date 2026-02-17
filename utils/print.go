package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mmoehabb/luci/types"
)

// PrintHeader displays the application logo, title, and description from the
// provided configuration. It uses colored output to make the header visually
// distinctive and wraps the description text for better readability.
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

// PrintUsage prints the complete usage information, including the header and
// all available actions from the shell configuration. It iterates through
// each action in the configuration and displays them in a formatted manner.
func PrintUsage(c types.Config) {
	PrintHeader(c)
	shell := *GetShellConfig(c)
	for action := range shell {
		PrintActionWithInputs(shell, []string{action}, 0)
	}
}

// PrintActionWithInputs resolves an action from the configuration using the
// provided inputs and prints it. It returns an error if the action cannot be
// found, otherwise nil on successful printing.
func PrintActionWithInputs(c map[string]any, inputs []string, level int) error {
	action := Dig(c, inputs)
	if action == nil {
		return errors.New("Action couldn't be found!")
	}
	PrintAction(action, inputs, level)
	return nil
}

// PrintAction prints an action in a formatted way based on its type. It handles
// AnnotatedAction, map[string]any, []string, and string types, applying
// appropriate colors and indentation to display hierarchical action structures.
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

// PrintCommand displays the command that is about to be executed with a
// highlighted green background and white text, making it visually distinct
// in the terminal output.
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
