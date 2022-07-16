package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	ch := make(chan int, 10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	go productor1(ch, ctx)
	go consumer1(ch)

	select {
	case <-ctx.Done():
		fmt.Println("main over")
		//等待其他go也结束
		time.Sleep(time.Second * 5)
		return
	}

}

func productor1(ch chan<- int, ctx context.Context) {
	i := 0
	for range time.Tick(time.Second) {
		select {
		case <-ctx.Done():
			fmt.Println("productor is over")
			close(ch)
			return
		default:
			fmt.Println("product to chan ", i)
			ch <- i
			i = i + 1
		}
	}
}

func consumer1(ch <-chan int) {
	for range time.Tick(time.Second) {
		if value, ok := <-ch; ok == true {
			fmt.Println("consumer from chan ", value)
		} else {
			fmt.Println("consumer is over ")
			return
		}
	}
}
