package beacon

type Message struct {
	Payload []byte
}

type RoutedMessage struct {
	Message
	Topic *TopicMatch
}

func (m *RoutedMessage) GetTopicParam(param string) string {
	return m.Topic.Params()[param]
}
