package goutils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofreego/goutils/constants"
	"github.com/gofreego/goutils/logger"
	"github.com/google/uuid"
)

func GetHTTPRouter(mode string) *gin.Engine {
	gin.SetMode(mode)
	engine := gin.New()
	engine.Use(RequestTimeMiddleware)
	engine.Use(RequestIDMiddleware)
	engine.Use(OptionRequestMiddleware)
	return engine
}

func OptionRequestMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
		return
	}
	c.Next()
}

func RequestTimeMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()
	field := &logger.Fields{}
	field.AddField("totalTime", time.Since(start).Milliseconds())
	logger.Infof(c, fmt.Sprintf("%s %s request end ", c.Request.Method, c.Request.URL.RequestURI()), field)

}

func RequestIDMiddleware(c *gin.Context) {
	// Set the request ID in the context
	if c.GetHeader(constants.X_REQUEST_ID) != "" {
		c.Set(constants.X_REQUEST_ID, uuid.New().String())
	}
	// Pass control to the next middleware or route handler
	c.Next()
}
