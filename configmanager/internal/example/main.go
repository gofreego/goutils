package example

import (
	"context"

	"github.com/gofreego/goutils/configmanager"
)

type Repo struct {
}

func (r *Repo) GetConfig(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (r *Repo) SaveConfig(ctx context.Context, key string, value string) error {
	return nil
}

func main() {
	configmanager.New(ctx, &Repo{})
}
