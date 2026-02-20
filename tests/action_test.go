package tests

import (
	"github.com/mmoehabb/luci/types"
	"github.com/mmoehabb/luci/utils"
	"reflect"
	"testing"
)

func TestDig(t *testing.T) {
	tests := []struct {
		name     string
		action   any
		inputs   []string
		expected any
		index    int
	}{
		{
			name: "find nested action in ShellConfig",
			action: types.ShellConfig{
				"deploy": types.ShellConfig{
					"prod": "echo deploying to production",
				},
			},
			inputs:   []string{"deploy", "prod"},
			expected: "echo deploying to production",
			index:    2,
		},
		{
			name: "find action in AnnotatedAction",
			action: types.AnnotatedAction{
				Value: types.ShellConfig{
					"build": "echo building",
				},
			},
			inputs:   []string{"build"},
			expected: "echo building",
			index:    1,
		},
		{
			name: "find action in nested AnnotatedAction",
			action: types.AnnotatedAction{
				Value: types.AnnotatedAction{
					Value: "echo nested action",
				},
			},
			inputs:   []string{"nested"},
			expected: "echo nested action",
			index:    1,
		},
		{
			name: "find action in map[string]any",
			action: map[string]any{
				"value": types.ShellConfig{
					"test": "echo testing",
				},
			},
			inputs:   []string{"test"},
			expected: "echo testing",
			index:    1,
		},
		{
			name: "return nil when no matching action",
			action: types.ShellConfig{
				"build": "echo building",
			},
			inputs:   []string{"deploy"},
			expected: nil,
			index:    1,
		},
		{
			name:     "handle empty inputs",
			action:   "echo hello",
			inputs:   []string{},
			expected: "echo hello",
			index:    0,
		},
		{
			name:     "handle empty action",
			action:   nil,
			inputs:   []string{"test"},
			expected: nil,
			index:    0,
		},
		{
			name: "handle complex nested structure",
			action: types.ShellConfig{
				"project": types.AnnotatedAction{
					Value: map[string]any{
						"value": types.ShellConfig{
							"setup": "echo setting up",
						},
					},
				},
			},
			inputs:   []string{"project", "setup"},
			expected: "echo setting up",
			index:    2,
		},
		{
			name: "handle AnnotatedAction with string value",
			action: types.AnnotatedAction{
				Value: "echo direct string",
			},
			inputs:   []string{"any"},
			expected: "echo direct string",
			index:    1,
		},
		{
			name: "handle AnnotatedAction with []string value",
			action: types.AnnotatedAction{
				Value: []string{"echo", "array", "action"},
			},
			inputs:   []string{"any"},
			expected: []string{"echo", "array", "action"},
			index:    1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, index := utils.Dig(tt.action, tt.inputs)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Dig() = %v, want %v", result, tt.expected)
			}

			if index != tt.index {
				t.Errorf("Dig() index = %d, want %d", index, tt.index)
			}
		})
	}
}

func TestMapToAnnotatedAction(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]any
		expected types.AnnotatedAction
	}{
		{
			name: "convert map with all fields",
			input: map[string]any{
				"title":       "Test Title",
				"description": "Test Description",
				"value":       "echo test",
			},
			expected: types.AnnotatedAction{
				Title:       "Test Title",
				Description: "Test Description",
				Value:       "echo test",
			},
		},
		{
			name: "convert map with missing optional fields",
			input: map[string]any{
				"value": "echo test",
			},
			expected: types.AnnotatedAction{
				Title:       "",
				Description: "",
				Value:       "echo test",
			},
		},
		{
			name: "convert map with only title",
			input: map[string]any{
				"title": "Only Title",
				"value": "echo test",
			},
			expected: types.AnnotatedAction{
				Title:       "Only Title",
				Description: "",
				Value:       "echo test",
			},
		},
		{
			name: "convert map with only description",
			input: map[string]any{
				"description": "Only Description",
				"value":       "echo test",
			},
			expected: types.AnnotatedAction{
				Title:       "",
				Description: "Only Description",
				Value:       "echo test",
			},
		},
		{
			name: "convert map with []string value",
			input: map[string]any{
				"value": []string{"echo", "test"},
			},
			expected: types.AnnotatedAction{
				Title:       "",
				Description: "",
				Value:       []string{"echo", "test"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := utils.MapToAnnotatedAction(tt.input)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("MapToAnnotatedAction() = %v, want %v", result, tt.expected)
			}
		})
	}
}
