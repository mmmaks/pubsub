package main

import (
	"context"
	pubsub "personal.com/m"
	"personal.com/m/broker"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	pubsubService := initialize(ctx)

	runCases(ctx, pubsubService)

	select {
	case <-ctx.Done():
	}
}

func runCases(ctx context.Context, pubsubService pubsub.Pubsub) {

	err := pubsubService.CreateTopic("topic1")
	if err != nil {
		panic(err)
	}
	err = pubsubService.CreateTopic("topic2")
	if err != nil {
		panic(err)
	}

	err = pubsubService.Subscribe(ctx, "sub1", "topic1", func(data string) error {
		println("sub1 received: ", data)
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = pubsubService.Subscribe(ctx, "sub2", "topic2", func(data string) error {
		println("sub2 received: ", data)
		return nil
	})
	if err != nil {
		panic(err)
	}

	err = pubsubService.Publish(ctx, "topic4", "hello1")
	if err != nil {
		println(err.Error())
	}
	err = pubsubService.Publish(ctx, "topic1", "world1")
	if err != nil {
		println(err.Error())
	}
	err = pubsubService.Publish(ctx, "topic2", "hello2")
	if err != nil {
		println(err.Error())
	}
	err = pubsubService.Publish(ctx, "topic2", "world2")
	if err != nil {
		println(err.Error())
	}
}

func initialize(ctx context.Context) pubsub.Pubsub {

	b := broker.NewBroker(ctx)

	p := pubsub.NewPubsub(b)
	return p
}
