package manifest

import (
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
)

type Data struct {
	Description  string
	Timeout      int
	Options      string
	RequiredVars []string
	AssetTypes   []string
}

// Read reads a manifest file.
func Read(path string) (Data, error) {
	d := Data{}
	m, err := toml.DecodeFile(path, &d)
	if err != nil {
		return Data{}, err
	}

	if m.IsDefined("Options") {
		dummy := make(map[string]interface{})
		err = json.Unmarshal([]byte(d.Options), &dummy)
		if err != nil {
			err = fmt.Errorf("options field is not a valid json string: %v", err)
			return Data{}, err
		}
	}
	return d, nil
}
