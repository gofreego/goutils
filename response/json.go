package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/customerrors"
)

type Response struct {
	Message *string `json:"message,omitempty"`
	Error   *string `json:"error,omitempty"`
	Data    any     `json:"data,omitempty"`
}

func WriteError(ctx *gin.Context, err error) {
	errStr := "something went wrong"
	if customErr, ok := err.(*customerrors.Error); ok {
		errStr = customErr.Error()
		ctx.JSON(customErr.Code, &Response{Error: &errStr})
		return
	}
	ctx.JSON(http.StatusInternalServerError, &Response{Error: &errStr})
}

func WriteSuccess(ctx *gin.Context, data any) {
	if msg, ok := data.(string); ok {
		ctx.JSON(http.StatusOK, &Response{Message: &msg})
		return
	}
	ctx.JSON(http.StatusOK, &Response{Data: data})
}
