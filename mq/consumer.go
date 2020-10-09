package mq

type Consumer struct {
	client *Client
	topic  string
	sub <-chan interface{}
}

func NewConsumer(c *Client, topic string) (*Consumer, error) {
	ch, err := c.Subscribe(topic)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		client: c,
		topic:  topic,
		sub: ch,
	}, nil
}

func (c *Consumer) Consume(f func(msg interface{}) error) error {
	for {
		msg := <-c.sub
		err := f(msg)
		if err != nil {
			return err
		}
	}
}
