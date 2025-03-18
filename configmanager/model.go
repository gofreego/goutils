package configmanager

import (
	"github.com/gofreego/goutils/customerrors"
)

type ConfigObject struct {
	Name        string         `json:"name"`
	Type        ConfigType     `json:"type"`
	Description string         `json:"description"`
	Required    bool           `json:"required"`
	Choices     []string       `json:"choices,omitempty"`
	Value       any            `json:"value"`
	Childrens   []ConfigObject `json:"children"`
}

func (co ConfigObject) Validate() error {

	if co.Required && co.Value == nil {
		return customerrors.BAD_REQUEST_ERROR("config %s is required, please pass the value", co.Name)
	}

	if co.Value != nil {
		switch co.Type {
		case CONFIG_TYPE_STRING, CONFIG_TYPE_BIG_TEXT:
			if _, ok := co.Value.(string); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: string", co.Name, co.Value)
			}
		case CONFIG_TYPE_NUMBER:
			if _, ok := co.Value.(float64); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: number", co.Name, co.Value)
			}
		case CONFIG_TYPE_BOOLEAN:
			if _, ok := co.Value.(bool); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: boolean", co.Name, co.Value)
			}
		case CONFIG_TYPE_JSON:
			if _, ok := co.Value.(map[string]any); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: json", co.Name, co.Value)
			}
		case CONFIG_TYPE_CHOICE:
			if _, ok := co.Value.(string); !ok {
				return customerrors.BAD_REQUEST_ERROR("config %s has invalid value type %T, Expect: string", co.Name, co.Value)
			}
		default:
			return customerrors.BAD_REQUEST_ERROR("config %s has invalid type %s", co.Name, co.Type)
		}
	}
	for _, child := range co.Childrens {
		if err := child.Validate(); err != nil {
			return err
		}
	}

	return nil
}
