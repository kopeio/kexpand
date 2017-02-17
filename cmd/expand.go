package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"encoding/base64"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

type ExpandCmd struct {
	cobraCommand *cobra.Command

	SourceFiles []string
	Values      []string

	IgnoreMissingFiles bool
	IgnoreMissingKeys  bool
}

var expandCmd = ExpandCmd{
	cobraCommand: &cobra.Command{
		Use:   "expand",
		Short: "Expand a template",
	},
}

func init() {
	cmd := expandCmd.cobraCommand
	rootCommand.cobraCommand.AddCommand(cmd)

	cmd.Flags().StringSliceVarP(&expandCmd.SourceFiles, "file", "f", nil, "files containing values to substitute")
	cmd.Flags().StringSliceVarP(&expandCmd.Values, "value", "k", nil, "key=value pairs to substitute")
	cmd.Flags().BoolVarP(&expandCmd.IgnoreMissingFiles, "ignore-missing-files", "i", false, "ignore source files that are not found")
	cmd.Flags().BoolVar(&expandCmd.IgnoreMissingKeys, "ignore-missing-keys", false, "ignore missing value keys that are not found")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		err := expandCmd.Run(args)
		if err != nil {
			glog.Exitf("%v", err)
		}
	}
}

func (c *ExpandCmd) Run(args []string) error {
	values, err := c.parseValues()
	if err != nil {
		return err
	}

	for k, v := range values {
		glog.V(2).Infof("\t%q=%q", k, v)
	}

	var src []byte
	if len(args) == 0 {
		src, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("error reading from stdin: %v", err)
		}
	} else if len(args) == 1 {
		src, err = ioutil.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("error reading file %q: %v", args[0], err)
		}
	} else {
		return fmt.Errorf("expected exactly one argument, a path to a file to expand")
	}

	expanded, err := c.DoExpand(src, values)

	_, err = os.Stdout.Write(expanded)
	if err != nil {
		return fmt.Errorf("error writing to stdout: %v", err)
	}

	return nil
}

func (c *ExpandCmd) parseValues() (map[string]interface{}, error) {
	values := make(map[string]interface{})

	for _, f := range c.SourceFiles {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			if c.IgnoreMissingFiles && os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Skipping missing file %q\n", f)
				continue
			}
			return nil, fmt.Errorf("error reading file %q: %v", f, err)
		}

		data, err := parseYamlSource(f, b)
		if err != nil {
			return nil, err
		}

		for k, v := range data {
			values[k] = v
		}
	}

	for _, v := range c.Values {
		tokens := strings.SplitN(v, "=", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("Unexpected value %q, expected key=value", v)
		}
		values[tokens[0]] = tokens[1]
	}

	return values, nil
}

func parseYamlSource(source string, b []byte) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	if err := yaml.Unmarshal(b, &data); err != nil {
		return nil, fmt.Errorf("error parsing yaml file %q: %v", source, err)
	}
	return data, nil
}

func (c *ExpandCmd) DoExpand(src []byte, values map[string]interface{}) ([]byte, error) {
	expanded := src

	{
		var err error

		// All
		expr := `\$(\({1,2})([[:alnum:]_\.\-]+)(\|base64)?\){1,2}|(\{{2})([[:alnum:]_\.\-]+)(\|base64)?\}{2}`
		re := regexp.MustCompile(expr)
		expandFunction := func(match []byte) []byte {
			re := regexp.MustCompile(expr)

			matchStr := string(match[:])
			result := re.FindStringSubmatch(matchStr)

			if result[0] != matchStr {
				glog.Fatalf("Unexpected match: %q", matchStr)
			}

			if result[2] == "" && result[5] == "" {
				glog.Fatalf("No variable defined within: %q", matchStr)
			}

			key := result[2] + result[5]
			replacement := values[key]

			if replacement == nil {
				if c.IgnoreMissingKeys == false {
					err = fmt.Errorf("Key not found: %q", key)
				}
				return match
			}

			if (result[3] + result[6]) == "|base64" {
				replacement = base64.StdEncoding.EncodeToString([]byte(replacement.(string)))
			}

			var s string
			delim := result[1] + result[4]
			switch len(delim) {
			case 1:
				s = fmt.Sprintf("\"%v\"", replacement)
			case 2:
				s = fmt.Sprintf("%v", replacement)
			default:
				glog.Fatalf("Unexpected delimiter %q count: %q", delim, len(delim))
			}

			return []byte(s)
		}

		expanded = re.ReplaceAllFunc(expanded, expandFunction)
		if err != nil {
			return nil, err
		}
	}

	return expanded, nil
}
