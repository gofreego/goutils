package configmanager

import "fmt"

var (
	ErrConfigNotFound = fmt.Errorf("config not found")
	ErrInvalidConfig  = fmt.Errorf("invalid config")
)
