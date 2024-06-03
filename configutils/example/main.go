package example

import (
	"context"
	"fmt"
	"time"

	"github.com/gofreego/goutils/configutils/consul"
	"github.com/gofreego/goutils/logger"
)

type conf struct {
	Name   string  `json:"name"`
	Int    int     `json:"int"`
	Float  float64 `json:"float"`
	Struct struct {
		Name string `json:"name"`
		Int  int    `json:"int"`
	} `json:"struct"`
}

func main() {
	config := consul.Config{
		Address:       "http://localhost:8500",
		Path:          "configs/test",
		RefreshInSecs: 2,
	}
	ctx := context.Background()
	reader, err := consul.NewConsulReader(ctx, &config)
	if err != nil {
		logger.Error(ctx, fmt.Sprintln("Error creating consul reader : ", err))
		return
	}
	var c conf
	err = reader.Read(&c)
	if err != nil {
		logger.Error(ctx, fmt.Sprintln("Error reading consul : ", err))
		return
	}
	for {
		time.Sleep(1 * time.Second)
		logger.Info(ctx, "%v", c)
	}
}
