package main

// import (
// 	"fmt"
// 	"sync"
// 	"time"
// )

// func main() {
// 	c := make(chan string)

// 	var wg sync.WaitGroup

// 	wg.Add(2)

// 	go func() {
// 		defer wg.Done()
// 		c <- "12345"
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		time.Sleep(1 * time.Second)
// 		fmt.Println("message:" + <-c)
// 	}()

// 	wg.Wait()
// }
