package eventqueue

import (
	"context"
	"fmt"

	"github.com/gofreego/goutils/eventqueue/kafka"
	"github.com/gofreego/goutils/eventqueue/models"
)

const (
	KAFKA = "kafka"
)

type EventQueue interface {
	//publish multiple messages
	Publish(ctx context.Context, events ...models.IMessage) error

	//consume single message
	Consume(ctx context.Context) (models.IMessage, error)
	ConsumeAndCommit(ctx context.Context) (models.IMessage, error)

	// consume many messages
	ConsumeMany(ctx context.Context) ([]models.IMessage, error)
	ConsumeManyAndCommit(ctx context.Context)

	// pass the consumed messages to commit
	Commit(ctx context.Context, events ...models.IMessage) error
}

type Config struct {
	Name  string
	Kafka kafka.Config
}

func NewEventQueue(ctx context.Context, cfg *Config) EventQueue {
	switch cfg.Name {
	case KAFKA:
		return kafka.NewEventQueue(ctx, &cfg.Kafka)
	}
	panic(fmt.Sprintf("invalid event queue name , provided : %s , expected : %s", cfg.Name, KAFKA))
}
