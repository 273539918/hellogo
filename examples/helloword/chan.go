package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	fmt.Println("-----------------------------demo 1------------------------------")
	ch := make(chan int, 1)
	for i := 0; i < 3; i++ {
		go func(i int) {
			fmt.Printf("hello from goroutine %d \n", i)
			//缓冲区满了会阻塞
			ch <- i
			fmt.Printf("value put chan %d \n", i)
		}(i)
	}
	for i := 0; i < 3; i++ {
		fmt.Printf("main loop %d \n", i)
		time.Sleep(time.Second * 1)
		//阻塞等待channel中的数据
		value := <-ch
		fmt.Printf("value from chan %d \n", value)
	}
	fmt.Println("-----------------------------demo 2------------------------------")
	ch2 := make(chan int, 10)
	go func() {
		defer close(ch2)
		for i := 0; i < 3; i++ {
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(10)
			fmt.Println("putting:", n)
			ch2 <- n
		}
	}()
	fmt.Println("main ")
	for v := range ch2 {
		fmt.Println("receiving:", v)
	}
	fmt.Println("-----------------------------demo 3------------------------------")

	ch3 := make(chan int, 10)
	chclose := make(chan bool)
	go productor(ch3, chclose)
	go consumer(ch3)
	time.Sleep(time.Second * 3)
	//chclose通道被关闭后，任意从chclose读取的操作都不会阻塞，直接返回 false
	close(chclose)

	time.Sleep(time.Second * 60)
}

func productor(productChan chan<- int, closeFlag <-chan bool) {
	i := 0
	for {
		select {
		case <-closeFlag:
			fmt.Println("product is close:")
			close(productChan)
			return
		default:
			fmt.Println("product :", i)
			productChan <- i
			i = i + 1
			//time.Sleep(time.Second)
		}

	}
}

func consumer(consumerChan <-chan int) {
	for {
		if num, ok := <-consumerChan; ok == true {
			fmt.Println("consumer : ", num)
		} else {
			fmt.Println("consumer is close ,but if continue to get value,will return :", <-consumerChan)
			return
		}
		time.Sleep(time.Second)
	}

}
