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

var successStatusCodesMap = map[int]bool{
	http.StatusOK:        true,
	http.StatusCreated:   true,
	http.StatusAccepted:  true,
	http.StatusNoContent: true,
}

func WriteSuccess(ctx *gin.Context, data any, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		if _, ok := successStatusCodesMap[statusCode[0]]; ok {
			code = statusCode[0]
		}
	}

	if msg, ok := data.(string); ok {
		ctx.JSON(code, &Response{Message: &msg})
		return
	}
	ctx.JSON(code, &Response{Data: data})
}
