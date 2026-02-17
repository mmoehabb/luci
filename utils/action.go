package utils

import (
	"github.com/mmoehabb/luci/types"
)

// Dig recursively traverses the configuration structure based on the provided
// input keys. It navigates through ShellConfig, AnnotatedAction, or map[string]any
// types to find and return the action matching the given inputs. Returns nil if
// no matching action is found.
func Dig(action any, inputs []string) any {
	for i, input := range inputs {
		switch actTyped := action.(type) {
		case types.ShellConfig:
			action = actTyped[input]
			continue

		case types.AnnotatedAction:
			action = Dig(actTyped.Value, inputs[i:])
			continue

		case map[string]any:
			if actTyped["value"] != nil {
				action = Dig(MapToAnnotatedAction(actTyped), inputs[i:])
				continue
			}
			action = actTyped[input]
			continue
		}

		break
	}

	// In case the action is annotated one, then ensure that it's being return
	// as AnnotatedAction
	switch actTyped := action.(type) {
	case map[string]any:
		if actTyped["value"] != nil {
			return MapToAnnotatedAction(actTyped)
		}
	}

	return action
}

// MapToAnnotatedAction converts a generic map[string]any to an AnnotatedAction.
// It extracts the title, description, and value fields from the map and returns
// a properly typed AnnotatedAction struct. Fields that are not present default
// to empty strings.
func MapToAnnotatedAction(m map[string]any) types.AnnotatedAction {
	title := ""
	if m["title"] != nil {
		title = m["title"].(string)
	}

	description := ""
	if m["description"] != nil {
		description = m["description"].(string)
	}

	return types.AnnotatedAction{
		Title:       title,
		Description: description,
		Value:       m["value"],
	}
}
