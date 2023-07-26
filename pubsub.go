package pubsub

import (
	"context"
	"personal.com/m/broker"
)

type Pubsub interface {
	CreateTopic(topicName string) error
	Publish(ctx context.Context, topicName string, data string) error
	Subscribe(ctx context.Context, subName, topicName string, callback func(string) error) error
}

type pubsub struct {
	broker broker.Broker
}

func NewPubsub(broker broker.Broker) *pubsub {
	return &pubsub{
		broker: broker,
	}
}

func (p pubsub) CreateTopic(topicName string) error {
	return p.broker.CreateTopic(topicName)
}

func (p pubsub) Publish(ctx context.Context, topicName string, data string) error {
	return p.broker.Publish(ctx, topicName, data)
}

func (p pubsub) Subscribe(ctx context.Context, subName, topicName string, callback func(string) error) error {
	return p.broker.Subscribe(ctx, subName, topicName, callback)
}
