package logger

import (
	"context"
)

func RequestMiddleLayer(ctx context.Context, msg string, fields *Fields) (context.Context, string, *Fields) {
	vRc, rcOk := ctx.Value(RequestContextKey).(RequestContext)
	if !rcOk {
		return ctx, msg, fields
	}
	if vRc.RequestID != "" {
		fields.AddField(requestIDKey, vRc.RequestID)
	}
	if vRc.AppID != "" {
		fields.AddField(appIDKey, vRc.AppID)
	}
	if vRc.UserID != "" {
		fields.AddField(userIDKey, vRc.UserID)
	}
	if vRc.Method != "" {
		fields.AddField(callerKey, vRc.Method)
	}
	if vRc.URI != "" {
		fields.AddField(uriKey, vRc.URI)
	}
	if vRc.IP != "" {
		fields.AddField(ipKey, vRc.IP)
	}
	return ctx, msg, fields
}

type ContextKey string

const (
	RequestContextKey ContextKey = "RequestContext"
)

type RequestContext struct {
	RequestID string
	AppID     string
	UserID    string
	Method    string
	URI       string
	IP        string
}
