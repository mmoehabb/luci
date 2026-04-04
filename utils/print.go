package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
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

func PrintInteractiveUsage(c types.Config) {
	shell := *GetShellConfig(c)
	actions := CollectActions(shell)
	displayLevel(c, []string{}, actions)
}

func displayLevel(c types.Config, path []string, actions []ActionNode) {
	title := "Select an action"
	if len(path) > 0 {
		title = fmt.Sprintf("Select an action (%s)", strings.Join(path, " > "))
	}

	options := []huh.Option[string]{}
	if len(path) > 0 {
		options = append(options, huh.NewOption("⬅ Back", BackKey))
	}

	for _, action := range actions {
		description := action.Description
		if description == "" {
			description = action.Title
		}
		marker := "▸"
		if action.IsGroup {
			marker = "▸"
		} else {
			marker = "•"
		}
		label := fmt.Sprintf("%s %s", marker, description)
		options = append(options, huh.NewOption(label, action.Key))
	}

	var selected string
	selection := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(&selected)

	err := selection.Run()
	if err != nil {
		color.New(color.FgYellow).Println("\nCancelled")
		return
	}

	if selected == BackKey {
		if len(path) > 0 {
			parentPath := path[:len(path)-1]
			shell := *GetShellConfig(c)
			var parentActions []ActionNode
			if len(parentPath) == 0 {
				parentActions = CollectActions(shell)
			} else {
				parentConfig, _ := Dig(shell, parentPath)
				if parentConfig != nil {
					parentActions = CollectActions(parentConfig.(map[string]any))
				}
			}
			displayLevel(c, parentPath, parentActions)
		}
		return
	}

	newPath := append(path, selected)
	shell := *GetShellConfig(c)
	selectedAction, _ := Dig(shell, newPath)

	if selectedAction == nil {
		displayLevel(c, path, actions)
		return
	}

	switch a := selectedAction.(type) {
	case types.AnnotatedAction:
		if nested, ok := a.Value.(map[string]any); ok {
			if _, hasValue := nested["value"]; hasValue {
				Act(c, newPath)
			} else {
				children := CollectActions(nested)
				displayLevel(c, newPath, children)
			}
		} else {
			Act(c, newPath)
		}
	case map[string]any:
		if a["value"] != nil {
			Act(c, newPath)
		} else {
			children := CollectActions(a)
			displayLevel(c, newPath, children)
		}
	case string:
		Act(c, newPath)
	case []string:
		Act(c, newPath)
	default:
		displayLevel(c, path, actions)
	}
}

func PrintUsage(c types.Config) {
	PrintHeader(c)
	shell := *GetShellConfig(c)
	for action := range shell {
		PrintActionWithInputs(shell, []string{action}, 0)
	}
}

func PrintActionWithInputs(c map[string]any, inputs []string, level int) error {
	action, _ := Dig(c, inputs)
	if action == nil {
		return errors.New("Action couldn't be found!")
	}
	PrintAction(action, inputs, level)
	return nil
}

func PrintAction(action any, inputs []string, level int) {
	switch action := action.(type) {
	case types.AnnotatedAction:
		switch annVal := action.Value.(type) {
		case map[string]any:
			color.New(color.FgMagenta).Printf("%sluci %s\n", indent(level), strings.Join(inputs, " "))
			if action.Title != "" {
				color.Blue("%s** %s **", indent(level+1), action.Title)
			}
			if action.Description != "" {
				color.Black("%s> %s", indent(level+1), action.Description)
			}
			if annVal["value"] != nil {
				PrintAction(MapToAnnotatedAction(annVal), inputs, level+1)
				return
			}
			PrintAction(action.Value, inputs, level+1)

		case []string:
			color.New(color.FgWhite).Printf("%sluci %s\t", indent(level), strings.Join(inputs, " "))
			if action.Title != "" {
				color.New(color.Faint).Printf("%s\t", action.Title)
			}
			color.New(color.Faint).Printf("%s\n", action.Description)

		case string:
			color.New(color.FgWhite).Printf("%sluci %s\t", indent(level), strings.Join(inputs, " "))
			if action.Title != "" {
				color.New(color.Faint).Printf("%s\t", action.Title)
			}
			color.New(color.Faint).Printf("%s\n", action.Description)

		default:
			PrintAction(action.Value, inputs, level+1)
		}

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

func PrintCommand(cmd string, args []string) {
	color.New(color.BgBlack, color.FgHiWhite).Printf("+ %s %s", cmd, strings.Join(args, " "))
	fmt.Println()
}

func indent(count int) string {
	var res strings.Builder
	for range count {
		res.WriteString("\t")
	}
	return res.String()
}
