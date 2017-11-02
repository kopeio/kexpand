package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/kopeio/kexpand/pkg/source"
	"github.com/spf13/cobra"
)

type ExpandCmd struct {
	cobraCommand *cobra.Command

	SourceFiles []string
	BlobTrees   []string
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
	cmd.Flags().StringSliceVarP(&expandCmd.BlobTrees, "tree", "t", nil, "directory tree of files")
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
	} else {
		for argi, filepattern := range args {
			files, err := filepath.Glob(filepattern)
			if err != nil {
				return fmt.Errorf("error bad pattern %q: %v", filepattern, err)
			}
			for filei, file := range files {
				filesrc, err := ioutil.ReadFile(file)
				if err != nil {
					return fmt.Errorf("error reading file %q: %v", file, err)
				}
				if argi == 0 && filei == 0 {
					src = filesrc
				} else {
					filesrc = append([]byte("\n---\n"), filesrc...)
					src = append(src, filesrc...)
				}
			}
		}
	}

	expanded, err := c.DoExpand(src, values)
	if err != nil {
		return fmt.Errorf("error expanding template: %v", err)
	}

	_, err = os.Stdout.Write(expanded)
	if err != nil {
		return fmt.Errorf("error writing to stdout: %v", err)
	}

	return nil
}

func (c *ExpandCmd) parseValues() (map[string]interface{}, error) {
	values := make(map[string]interface{})

	for k, v := range c.defaultValues() {
		values[k] = v
	}

	for _, f := range c.BlobTrees {
		treeSource := &source.FiletreeSource{}
		data, err := treeSource.Build(f, "file.")
		if err != nil {
			return nil, err
		}

		for k, v := range data {
			values[k] = v
		}
	}

	for _, f := range c.SourceFiles {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			if c.IgnoreMissingFiles && os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Skipping missing file %q\n", f)
				continue
			}
			return nil, fmt.Errorf("error reading file %q: %v", f, err)
		}

		var data map[string]interface{}
		name := strings.ToLower(path.Base(f))
		if strings.HasSuffix(name, ".tfstate") {
			tfParser := &source.TerraformSource{}
			data, err = tfParser.Parse(b)
			if err != nil {
				return nil, err
			}
		} else {
			data, err = parseYamlSource(f, b)
			if err != nil {
				return nil, err
			}
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
		expr := `\$(\({1,2})([[:alnum:]_\.\-]+)(\|[[:alnum:]]+)?\){1,2}|(\{{2})([[:alnum:]_\.\-]+)(\|[[:alnum:]]+)?\}{2}`
		re := regexp.MustCompile(expr)
		expandFunction := func(match []byte) []byte {
			var replacement interface{}

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

			if value, ok := os.LookupEnv(key); ok {
				replacement = value
			} else {
				replacement = values[key]
				if replacement == nil {
					if c.IgnoreMissingKeys == false {
						err = fmt.Errorf("Key not found: %q", key)
					}
					return match
				}
			}

			pipeFunction := result[3] + result[6]
			if pipeFunction != "" {
				if pipeFunction == "|base64" {
					b, ok := replacement.([]byte)
					if !ok {
						b = []byte(replacement.(string))
					}
					replacement = base64.StdEncoding.EncodeToString(b)
				} else if pipeFunction == "|yaml" {
					b, ok := replacement.(string)
					if !ok {
						b = string(replacement.([]byte))
					}
					data, err := json.Marshal(b)
					if err != nil {
						glog.Fatalf("error converting to JSON/YAML: %v", err)
					}
					replacement = string(data)
				} else {
					glog.Fatalf("Unknown pipe function: %q", pipeFunction)
				}
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

func (c *ExpandCmd) defaultValues() map[string]interface{} {
	defaults := make(map[string]interface{})

	t := time.Now()
	defaults["ansicdate"] = t.Format(time.ANSIC)
	defaults["ansicdateutc"] = t.UTC().Format(time.ANSIC)
	defaults["rubydate"] = t.Format(time.RubyDate)
	defaults["rubydateutc"] = t.UTC().Format(time.RubyDate)
	defaults["unixdate"] = t.Format(time.UnixDate)
	defaults["unixdateutc"] = t.UTC().Format(time.UnixDate)
	defaults["unixtime"] = t.Unix()
	defaults["unixtimeutc"] = t.UTC().Unix()

	wd, err := os.Getwd()
	if err == nil {
		defaults["basename"] = filepath.Base(wd)
		defaults["dirname"] = filepath.Dir(wd)
	}

	gitsha, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err == nil {
		defaults["gitsha"] = strings.TrimSpace(string(gitsha[:]))
		defaults["gitshashort"] = string(gitsha[:7])
	}

	return defaults
}
