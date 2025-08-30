package response

import "github.com/gin-gonic/gin"

type HandleFunc func(ctx *gin.Context) (interface{}, error)

func WithResponseMiddleware(f HandleFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		data, err := f(ctx)
		if err != nil {
			WriteError(ctx, err)
			return
		}
		WriteSuccess(ctx, data)
	}
}
