package configmanager

type ConfigObject struct {
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Required    bool           `json:"required"`
	Choices     []string       `json:"choices,omitempty"`
	Value       any            `json:"value"`
	Address     any            `json:"-"`
	Childrens   []ConfigObject `json:"children"`
}

func (co ConfigObject) Validate() error {
	return nil
}
