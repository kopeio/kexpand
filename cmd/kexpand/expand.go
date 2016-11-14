package main

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
	"github.com/kopeio/kexpand/pkg/expand"
)

type ExpandCmd struct {
	cobraCommand *cobra.Command

	SourceFiles []string
	Values      []string

	IgnoreMissingFiles bool
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

	expanded, err := expand.DoExpand(src, values)
	if err != nil {
		return err
	}

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

		data := make(map[string]interface{})
		if err := yaml.Unmarshal(b, &data); err != nil {
			return nil, fmt.Errorf("error parsing yaml file %q: %v", f, err)
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
