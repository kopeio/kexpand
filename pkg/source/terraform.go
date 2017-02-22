package source

import (
	"encoding/json"
	"fmt"
	"strings"
)

type TerraformSource struct {
}

type tfState struct {
	Version          int    `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Serial           int    `json:"serial"`
	Lineage          string `json:"lineage"`

	Modules []*tfStateModule `json:"modules"`
}

type tfStateModule struct {
	Path      []string                    `json:"path"`
	Resources map[string]*tfStateResource `json:"resources"`
	Outputs   map[string]*tfStateOutput   `json:"outputs"`

	DependsOn []string `json:"depends_on"`
}

type tfStateResource struct {
	Type      string         `json:"type"`
	DependsOn []string       `json:"depends_on"`
	Primary   *tfStateObject `json:"primary"`
	//"deposed": [],
	//"provider": ""
}

type tfStateObject struct {
	Id         string                 `json:"id"`
	Attributes map[string]interface{} `json:"attributes"`
}

type tfStateOutput struct {
	Sensitive bool        `json:"sensitive"`
	Type      string      `json:"type"`
	Value     interface{} `json:"value"`
}

func (t *TerraformSource) Parse(src []byte) (map[string]interface{}, error) {
	tf := &tfState{}
	err := json.Unmarshal(src, tf)
	if err != nil {
		return nil, fmt.Errorf("error parsing tfstate file: %v", err)
	}

	values := make(map[string]interface{})

	for _, module := range tf.Modules {
		path := []string{"tf"}
		for i, p := range module.Path {
			if i == 0 && p == "root" {
				continue
			}
			path = append(path, p)
		}
		prefix := strings.Join(path, ".")
		if prefix != "" {
			prefix += "."
		}

		for key, output := range module.Outputs {
			values[prefix+key] = output.Value
		}

		for key, resource := range module.Resources {
			if resource.Primary == nil {
				continue
			}
			for k, v := range resource.Primary.Attributes {
				values[prefix+key+"."+k] = v
			}
		}
	}

	return values, nil
}
