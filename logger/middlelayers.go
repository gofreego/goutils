package logger

import (
	"context"
	"time"
)

func RequestMiddleLayer(ctx context.Context, msg string, fields *Fields) (context.Context, string, *Fields) {
	vRc, rcOk := ctx.Value(RequestContextKey).(RequestContext)
	if !rcOk {
		return ctx, msg, fields
	}
	fields.AddField(requestIDKey, vRc.RequestID)
	fields.AddField(appIDKey, vRc.ClientAppID)
	fields.AddField(userIDKey, vRc.UserID)
	fields.AddField(uriKey, vRc.URI)
	fields.AddField(ipKey, vRc.IP)
	return ctx, msg, fields
}

const (
	RequestContextKey string = "RequestContext"
)

type RequestContext struct {
	SessionID     string
	RequestID     string
	ClientAppID   string
	UserID        string
	TransactionID string
	Method        string
	URI           string
	APIStartTime  time.Time
	IP            string
}
