package topics

import "context"

type Topic interface {
	Publish(ctx context.Context, data string) error
	Data() chan string
}

type topic struct {
	name string
	data chan string
}

func NewTopic(name string, size int) *topic {
	return &topic{
		name: name,
		data: make(chan string, size),
	}
}

func (t topic) Publish(ctx context.Context, data string) error {
	t.data <- data
	return nil
}

func (t topic) Data() chan string {
	return t.data
}
