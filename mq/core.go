package mq

import (
	"errors"
	"sync"
	"time"
)

//Broker 接口
type Broker interface {
	publish(topic string, msg interface{}) error
	subscribe(topic string) (<-chan interface{}, error)
	unsubscribe(topic string, sub <-chan interface{}) error
	close()
	broadcast(msg interface{}, subscribers []chan interface{})
	setConditions(capacity int)
}

//BrokerImpl 结构定义
type BrokerImpl struct {
	//关闭消息队列
	exit chan bool
	//消息队列的容量
	capacity int

	//key-topic,value-对应订阅者的channel
	topics map[string][]chan interface{}
	//读写锁
	sync.RWMutex
}

// NewBroker 创建新的Broker示例
func NewBroker() *BrokerImpl {
	return &BrokerImpl{
		exit:   make(chan bool),
		topics: make(map[string][]chan interface{}),
	}
}

func (b *BrokerImpl) setConditions(capacity int) {
	b.capacity = capacity
}

func (b *BrokerImpl) close() {
	select {
	case <-b.exit:
		return
	default:
		close(b.exit)
		b.Lock()
		b.topics = make(map[string][]chan interface{})
		b.Unlock()
	}
}

func (b *BrokerImpl) publish(topic string, msg interface{}) error {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.Lock()
	subscribers, ok := b.topics[topic]
	b.Unlock()
	if !ok {
		return nil
	}

	b.broadcast(msg, subscribers)
	return nil
}

func (b *BrokerImpl) broadcast(msg interface{}, subscribers []chan interface{}) {
	count := len(subscribers)
	concurrency := 1

	switch {
	case count > 1000:
		concurrency = 3
	case count > 100:
		concurrency = 2
	default:
		concurrency = 1
	}

	pub := func(start int) {
		for i := start; i < count; i += concurrency {
			select {
			case subscribers[i] <- msg:
			case <-time.After(time.Millisecond * 5):
			case <-b.exit:
				return
			}
		}
	}

	for j := 0; j < concurrency; j++ {
		go pub(j)
	}
}

func (b *BrokerImpl) subscribe(topic string) (<-chan interface{}, error) {
	select {
	case <-b.exit:
		return nil, errors.New("broker closed")
	default:
	}
	
	ch := make(chan interface{}, b.capacity)
	b.Lock()
	b.topics[topic] = append(b.topics[topic], ch)
	b.Unlock()
	return ch, nil;
}

func (b *BrokerImpl) unsubscribe(topic string, sub <-chan interface{}) (error) {
	select {
	case <-b.exit:
		return errors.New("broker closed")
	default:
	}

	b.RLock()
	subs, ok := b.topics[topic]
	b.RUnlock()
	if !ok {
		return nil
	}

	var newSubs []chan interface{}
	for _, subscribe := range subs {
		if (subscribe == sub) {
			continue
		}
		newSubs = append(newSubs, subscribe)
	}
	b.Lock()
	b.topics[topic] = newSubs
	b.Unlock()
	return nil
}

