package cmd

import "testing"

func Test_Find(t *testing.T) {
	type test struct {
		input    string
		expected string
		values   map[string]interface{}
	}
	grid := []test{
		{
			input:    "hello",
			expected: "hello",
		},
		{
			input:    "hello $((who))",
			expected: "hello world",
			values: map[string]interface{}{
				"who": "world",
			},
		},
		{
			input:    "hello {{who}}",
			expected: "hello world",
			values: map[string]interface{}{
				"who": "world",
			},
		},
		{
			input:    "hello $(who)",
			expected: "hello \"world\"",
			values: map[string]interface{}{
				"who": "world",
			},
		},
		{
			input:    "hello $((who|base64))",
			expected: "hello d29ybGQ=",
			values: map[string]interface{}{
				"who": "world",
			},
		},
	}

	for i, spec := range grid {
		c := &ExpandCmd{}
		actual, err := c.DoExpand([]byte(spec.input), spec.values)
		if err != nil {
			t.Errorf("doExpand unexpected error in test %d: %v", i, err)
		}

		if string(actual) != spec.expected {
			t.Errorf("unexpected expansion; expected=%q; actual=%q", spec.expected, string(actual))
		}
	}

}
