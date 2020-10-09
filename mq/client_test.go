package mq

import (
	"fmt"
	"sync"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient()
	client.setConditions(100)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		topic := fmt.Sprintf("mq test:%d", i)
		payload := fmt.Sprintf("shidongxuan666:%d", i)

		ch, err := client.Subscribe(topic)
		if err != nil {
			t.Fatal(err)
		}

		wg.Add(1)
		go func() {
			e := client.GetPayLoad(ch)
			if e != payload {
				t.Fatalf("%s expected %s but get %s", topic, payload, e)
			}
			if err := client.Unsubscribe(topic, ch); err != nil {
				t.Fatal(err)
			}
			wg.Done()
		}()

		if err := client.Publish(topic, payload); err != nil {
			t.Fatal(err)
		}
	}

	wg.Wait()
}