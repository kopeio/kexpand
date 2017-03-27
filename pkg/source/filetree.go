package source

import (
	"fmt"
	"io/ioutil"
	"path"
)

type FiletreeSource struct {
}

func (t *FiletreeSource) Build(basedir string, prefix string) (map[string]interface{}, error) {
	values := make(map[string]interface{})

	err := t.addBlobs(values, basedir, prefix)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (t *FiletreeSource) addBlobs(dest map[string]interface{}, dir string, keyPrefix string) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("error listing directory %q: %v", dir, err)
	}

	for _, f := range files {
		childKey := keyPrefix + f.Name()
		childPath := path.Join(dir, f.Name())
		if f.IsDir() {
			err := t.addBlobs(dest, childPath, childKey+".")
			if err != nil {
				return err
			}
		} else {
			contents, err := ioutil.ReadFile(childPath)
			if err != nil {
				return fmt.Errorf("error reading file %q: %v", childPath, err)
			}

			dest[childKey] = string(contents)
		}
	}

	return nil
}
