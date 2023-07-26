package exceptions

type TopicNotFound struct {
	TopicName string
}

func TopicNotFoundError(name string) error {
	return TopicNotFound{
		TopicName: name,
	}
}

func (e TopicNotFound) Error() string {
	return "topic not found: " + e.TopicName
}
