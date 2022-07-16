package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int)
	//一秒后触发一次 event
	timer := time.NewTimer(time.Second)
	select {
	case <-ch:
		fmt.Println("received from ch")
	case <-timer.C:
		fmt.Println("timeout waiting from channel ch")
	}

	//每秒执行一次
	for range time.Tick(time.Second) {
		fmt.Println("for while second trigger")
	}
}
