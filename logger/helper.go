package logger

import (
	"context"
)

func executeMiddleLayers(ctx context.Context, msg string, fields *Fields) (context.Context, string, *Fields) {

	for _, f := range middleLayers {
		ctx, msg, fields = f(ctx, msg, fields)
	}
	return ctx, msg, fields
}
