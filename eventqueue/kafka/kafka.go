package kafka

import (
	"context"

	"github.com/gofreego/goutils/eventqueue/models"
)

type Config struct {
	Brokers       []string
	Topic         string
	GroupID       string
	ConsumerGroup string
	Version       string
	IsPublisher   bool
}

type EventQueue struct {
}

// Commit implements eventqueue.EventQueue.
func (e *EventQueue) Commit(ctx context.Context, events ...models.IMessage) error {
	panic("unimplemented")
}

// Consume implements eventqueue.EventQueue.
func (e *EventQueue) Consume(ctx context.Context) (models.IMessage, error) {
	panic("unimplemented")
}

// ConsumeAndCommit implements eventqueue.EventQueue.
func (e *EventQueue) ConsumeAndCommit(ctx context.Context) (models.IMessage, error) {
	panic("unimplemented")
}

// ConsumeMany implements eventqueue.EventQueue.
func (e *EventQueue) ConsumeMany(ctx context.Context) ([]models.IMessage, error) {
	panic("unimplemented")
}

// ConsumeManyAndCommit implements eventqueue.EventQueue.
func (e *EventQueue) ConsumeManyAndCommit(ctx context.Context) {
	panic("unimplemented")
}

// Publish implements eventqueue.EventQueue.
func (e *EventQueue) Publish(ctx context.Context, events ...models.IMessage) error {
	panic("unimplemented")
}

func NewEventQueue(ctx context.Context, cfg *Config) *EventQueue {
	return &EventQueue{}
}
