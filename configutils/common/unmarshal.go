package common

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// Unmarshal unmarshals the data into the given config object.
// It supports json and yaml formats.
// If no format is provided, it defaults to yaml.
func Unmarshal(data []byte, conf any, cft ...ConfigFormatType) error {
	t := ConfigFormatYAML
	if len(cft) > 0 {
		t = cft[0]
	}

	switch t {
	case ConfigFormatJSON:
		return json.Unmarshal(data, conf)
	case ConfigFormatYAML:
		return yaml.Unmarshal(data, conf)
	}
	return ErrConfigFormatNotSupported
}

func Marshal(conf any, cft ...ConfigFormatType) ([]byte, error) {
	t := ConfigFormatYAML
	if len(cft) > 0 {
		t = cft[0]
	}

	switch t {
	case ConfigFormatJSON:
		return json.Marshal(conf)
	case ConfigFormatYAML:
		return yaml.Marshal(conf)
	}
	return nil, ErrConfigFormatNotSupported
}
