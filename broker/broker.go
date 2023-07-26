package broker

import (
	"context"
	"personal.com/m/exceptions"
	"personal.com/m/subscribers"
	"personal.com/m/topics"
	"sync"
)

type Broker interface {
	CreateTopic(topicName string) error
	Publish(ctx context.Context, topic string, data string) error
	Subscribe(ctx context.Context, subName, topic string, callback func(string) error) error
}

type broker struct {
	ctx         context.Context
	topics      map[string]topics.Topic
	subscribers map[string][]subscribers.Subscriber
	mu          sync.Mutex
}

func NewBroker(ctx context.Context) *broker {
	b := &broker{
		ctx:         ctx,
		topics:      make(map[string]topics.Topic),
		subscribers: make(map[string][]subscribers.Subscriber),
		mu:          sync.Mutex{},
	}
	return b
}

func (b *broker) CreateTopic(topicName string) error {

	defer b.synchronise()()

	if _, ok := b.topics[topicName]; ok {
		return nil
	}
	b.topics[topicName] = topics.NewTopic(topicName, 5)
	err := b.listenAndBroadcast(topicName)
	return err
}

func (b *broker) Publish(ctx context.Context, topic string, data string) error {

	if _, ok := b.topics[topic]; !ok {
		return exceptions.TopicNotFoundError(topic)
	}
	return b.topics[topic].Publish(ctx, data)
}

func (b *broker) Subscribe(ctx context.Context, subName, topic string, callback func(string) error) error {

	defer b.synchronise()()

	if _, ok := b.topics[topic]; !ok {
		return exceptions.TopicNotFoundError(topic)
	}
	s := subscribers.NewSubscriber(subName, topic, callback)
	b.subscribers[topic] = append(b.subscribers[topic], s)
	return s.Subscribe(ctx, topic, callback)
}

func (b *broker) listenAndBroadcast(topicName string) error {

	go func() {
		for {
			select {
			case msg := <-b.topics[topicName].Data():
				for _, sub := range b.subscribers[topicName] {
					sub.Msg() <- msg
				}
			case <-b.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (b *broker) synchronise() func() {

	b.mu.Lock()
	return b.mu.Unlock
}
