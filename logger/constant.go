package logger

const (
	// key constants
	timeKey      = "timeStamp"
	requestIDKey = "requestId"
	callerKey    = "caller"
	clientKey    = "client"
	userIDKey    = "userId"
	uriKey       = "uri"
	ipKey        = "ip"
	methodKey    = "method"
)

type Build string

const (
	// BuildProd : prod build
	BuildProd Build = "prod"
	// BuildDev : dev build
	BuildDev Build = "dev"
	// BuildTest : test build
)
