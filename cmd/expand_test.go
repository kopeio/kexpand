package cmd

import (
	"reflect"
	"testing"
)

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
		{
			input:    "hello $(who|base64)",
			expected: "hello \"d29ybGQ=\"",
			values: map[string]interface{}{
				"who": "world",
			},
		},
		{
			input:    "hello: $((who|yaml))",
			expected: "hello: \"hello\\nworld of yaml\"",
			values: map[string]interface{}{
				"who": "hello\nworld of yaml",
			},
		},
		{
			input:    "hello: $(who|yaml)",
			expected: "hello: \"\"hello\\nworld of yaml\"\"",
			values: map[string]interface{}{
				"who": "hello\nworld of yaml",
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
			t.Errorf("unexpected expansion of %q; expected=%q; actual=%q", spec.input, spec.expected, string(actual))
		}
	}
}

func TestParseYaml(t *testing.T) {
	grid := []struct {
		YAML     string
		Expected map[string]interface{}
	}{
		{
			YAML:     "",
			Expected: map[string]interface{}{},
		},
		{
			YAML:     "\n",
			Expected: map[string]interface{}{},
		},
		{
			YAML:     "# Nothing\n",
			Expected: map[string]interface{}{},
		},
		{
			YAML: "a: b",
			Expected: map[string]interface{}{
				"a": "b",
			},
		},
		{
			YAML: "# A comment\na: b\na2: 2\n",
			Expected: map[string]interface{}{
				"a":  "b",
				"a2": 2.0,
			},
		},
	}
	for _, g := range grid {
		actual, err := parseYamlSource("", []byte(g.YAML))
		if err != nil {
			t.Errorf("error parsing YAML %q: %v", g.YAML, err)
		}

		if !mapEquals(t, g.Expected, actual) {
			t.Errorf("Unexpected decoded value for %q.  Actual=%v, Expected=%v", g.YAML, actual, g.Expected)
		}

	}
}

func mapEquals(t *testing.T, l, r map[string]interface{}) bool {
	if len(l) != len(r) {
		return false
	}
	for k, vL := range l {
		vR := r[k]
		if !reflect.DeepEqual(vL, vR) {
			t.Logf("Not equals: %T %v vs %T %v", vL, vL, vR, vR)
			return false
		}
	}

	return true
}
