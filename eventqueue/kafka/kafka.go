package kafka

import (
	"context"
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

func NewEventQueue(ctx context.Context, cfg *Config) *EventQueue {
	return &EventQueue{}
}
