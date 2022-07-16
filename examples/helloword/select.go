package main

import (
	"fmt"
	"time"
)

func main() {

	ch1 := make(chan int)
	ch2 := make(chan int)

	go func(product1 chan<- int) {
		for {
			product1 <- 1
			time.Sleep(time.Second * 2)
		}
	}(ch1)

	go func(product1 chan<- int) {
		for {
			product1 <- 2
			time.Sleep(time.Second * 4)
		}
	}(ch2)

	go func() {
		for {
			select {
			case v := <-ch1:
				fmt.Println("v is ", v)
			case v := <-ch2:
				fmt.Println("v is ", v)
			default:
				fmt.Println("v is ", 0)
			}
			time.Sleep(time.Second)
		}
	}()

	time.Sleep(time.Second * 30)

}
