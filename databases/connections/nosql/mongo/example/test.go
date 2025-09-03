package main

import (
	"context"

	"github.com/gofreego/goutils/databases/connections/nosql/mongo"
)

func main() {
	conn, err := mongo.NewMongoConnection(context.Background(), &mongo.Config{
		Hosts:    "localhost:27017",
		Database: "abc",
		Username: "admin",
		Password: "XXXXX",
		Direct:   false,
	})
	if err != nil {
		panic(err)
	}
	defer conn.Disconnect(context.Background())
}
