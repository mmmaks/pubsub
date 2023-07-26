package subscribers

import "context"

type Subscriber interface {
	Subscribe(ctx context.Context, topicName string, callback func(string) error) error
	Msg() chan string
}

type subscriber struct {
	name      string
	topicName string
	exec      func(string) error
	msg       chan string
}

func NewSubscriber(name, topicName string, exec func(string) error) *subscriber {
	return &subscriber{
		name:      name,
		topicName: topicName,
		exec:      exec,
		msg:       make(chan string, 5),
	}
}

func (s subscriber) Subscribe(ctx context.Context, topicName string, callback func(string) error) error {

	go func() {
		for {
			select {
			case msg := <-s.msg:
				err := callback(msg)
				if err != nil {
					s.msg <- msg
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (s subscriber) Msg() chan string {
	return s.msg
}
