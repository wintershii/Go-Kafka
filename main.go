package main

import (
	"sync"
	"fmt"
	"main/mq"
	"time"
)

func main() {
	//OnceTopic()
	client := mq.NewClient()
	producer := mq.NewProducer(client)
	consumer, err := mq.NewConsumer(client, "1233")
	consumer2, err := mq.NewConsumer(client, "1233")
	if err != nil {
		panic("consumer init fialed")
	}
	go produce(producer)

	f := func(msg interface{}) error {
		fmt.Println("consume1 msg: %s", msg)
		return nil
	}

	f2 := func(msg interface{}) error {
		fmt.Println("consume2 msg: %s", msg)
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		consume(consumer, f)
	} ()
	go func() {
		defer wg.Done()
		consume(consumer2, f2)
	} ()
	wg.Wait()
}

func produce(producer *mq.Producer) {
	for i := 0; i < 100; i++ {
		producer.Send("1233", fmt.Sprintf("师东璇真帅:%d", i))
		time.Sleep(3 * time.Second)
	}
}

func consume(consumer *mq.Consumer, f func(msg interface{}) error) error {
	return consumer.Consume(f)
}


// func OnceTopic()  {
// 	client := mq.NewClient()
// 	client.SetConditions(10)
// 	ch, err := client.Subscribe("test-mq")
// 	if err != nil{
// 		fmt.Println("subscribe failed")
// 		return
// 	}
// 	go OncePub(client)
// 	OnceSub(ch, client)
// 	defer client.Close()
// }
   
// // 定时推送
// func OncePub(c *mq.Client) {
// 	t := time.NewTicker(10 * time.Second)
// 	defer t.Stop()
// 	for {
// 		select {
// 		case <- t.C:
// 			err := c.Publish("test-mq", "shidongxuan真帅")
// 	  		if err != nil {
// 	   			fmt.Println("pub message failed")
// 	  		}
// 		default:
   
// 	 	}
// 	}
// }
   
// // 接受订阅消息
// func OnceSub(msgs <-chan interface{}, c *mq.Client)  {
// 	for	{
// 		val := c.GetPayLoad(msgs)
// 		fmt.Printf("get message is %s\n",val)
// 	}
// }