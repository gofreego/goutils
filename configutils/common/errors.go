package common

import (
	"github.com/gofreego/goutils/customerrors"
)

var (
	ErrConfigFormatNotSupported = customerrors.BAD_REQUEST_ERROR("config format not supported, Expect one of json, yaml")
	ErrInvalidConfigReaderName  = customerrors.BAD_REQUEST_ERROR("invalid config reader name, Expect one of consul, zookeeper, database, file")
	ErrInvalidConfig            = customerrors.BAD_REQUEST_ERROR("invalid config")
)
