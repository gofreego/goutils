package logger

const (
	// key constants
	timeKey      = "TimeStamp"
	requestIDKey = "RequestID"
	callerKey    = "caller"
	appIDKey     = "AppId"
	userIDKey    = "UserID"
	uriKey       = "URI"
	ipKey        = "IP"
	methodKey    = "Method"
)

type Build string

const (
	// BuildProd : prod build
	BuildProd Build = "prod"
	// BuildDev : dev build
	BuildDev Build = "dev"
	// BuildTest : test build
)
