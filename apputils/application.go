package apputils

import "context"

type Application interface {
	Name() string
	Run(ctx context.Context) error
	Shutdown(ctx context.Context)
}
