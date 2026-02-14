package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/mmoehabb/luci/types"
)

func PrintHeader(c types.Config) {
	color.HiGreen(`*** %s ***`, c.Title)
	color.HiWhite(`> %s`, c.Description)
	color.Yellow("\nUsage:\n\n")
}

func PrintUsage(c types.Config) {
	PrintHeader(c)
	for action := range *GetShellConfig(c) {
		PrintAction(c, []string{action}, 0)
	}
}

func PrintAction(c types.Config, inputs []string, level int) {
	action := Dig(c, inputs)
	if action == nil {
		PrintUsage(c)
		return
	}

	switch Dig(c, inputs).(type) {
	case types.AnnotatedAction:
		color.New(color.FgMagenta).Printf("%sluci %s\n", indent(level), strings.Join(inputs, " "))
		annAct := action.(types.AnnotatedAction)
		if annAct.Title != "" {
			color.Blue("%s** %s **", indent(level+1), annAct.Title)
		}
		if annAct.Description != "" {
			color.Black("%s> %s", indent(level+1), annAct.Description)
		}
		if reflect.ValueOf(annAct.Value).Kind() == reflect.Map {
			for key := range annAct.Value.(map[string]any) {
				PrintAction(c, append(inputs, key), level+1)
			}
		}
		return

	case map[string]any:
		m := action.(map[string]any)
		for key := range m {
			PrintAction(c, append(inputs, key), level)
		}

	case []string:
		fmt.Printf("%sluci %s\n", indent(level), strings.Join(inputs, " "))

	case string:
		fmt.Printf("%sluci %s\n", indent(level), strings.Join(inputs, " "))
	}
}

func indent(count int) string {
	var res strings.Builder
	for range count {
		res.WriteString("\t")
	}
	return res.String()
}
