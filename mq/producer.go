package mq

type Producer struct {
	client *Client
}

func NewProducer(c *Client) *Producer {
	return &Producer{
		client : c,
	}
}

func (p *Producer) Send(topic string, msg interface{}) error {
	return p.client.Publish(topic, msg)
}

