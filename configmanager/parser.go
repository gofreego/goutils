package configmanager

import (
	"context"

	"github.com/gofreego/goutils/customerrors"
	"github.com/gofreego/goutils/datastructure"
)

var (
	// configTypesSet is a set of all valid config types.
	configTypesSet datastructure.Set[ConfigType] = datastructure.NewSet(
		CONFIG_TYPE_STRING,
		CONFIG_TYPE_INTEGER,
		CONFIG_TYPE_BOOLEAN,
		CONFIG_TYPE_JSON,
		CONFIG_TYPE_BIG_TEXT,
	)
)

func marshal(ctx context.Context, cfg config) (string, error) {
	return "", customerrors.New(0, "not implemented")
}

func unmarshal(ctx context.Context, value string, cfg config) error {
	return customerrors.New(0, "not implemented")
}
