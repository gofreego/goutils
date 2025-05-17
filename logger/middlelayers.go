package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
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
		fields.AddField(methodKey, vRc.Method)
	}
	if vRc.URI != "" {
		fields.AddField(uriKey, vRc.URI)
	}
	if vRc.IP != "" {
		fields.AddField(ipKey, vRc.IP)
	}
	return ctx, msg, fields
}

// Middleware to add Request ID
func WithRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		// set the request id in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, RequestContextKey, RequestContext{RequestID: requestID, URI: r.RequestURI, Method: r.Method, IP: r.RemoteAddr})
		r = r.WithContext(ctx) // update the request with the new context
		next.ServeHTTP(w, r)
	})
}

// request time middleware
func WithRequestTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// start time
		startTime := time.Now()
		next.ServeHTTP(w, r)
		Infow(r.Context(), "Request End", NewFields().AddField("totalTime", time.Since(startTime).Milliseconds()))
	})
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
