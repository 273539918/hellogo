package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	waitByTimeSleep()
	waitByChan()
	waitbyWG()
}

//通过time.sleep的等待让print输出完
func waitByTimeSleep() {
	fmt.Println("waitByTimeSleep")
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
		}(i)
	}
	time.Sleep(time.Second * 1)

}

//通过chan来让print输出完
func waitByChan() {
	fmt.Println("waitByChan")
	c := make(chan int, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			fmt.Println(i)
			c <- i
		}(i)
	}
	for i := 0; i < 10; i++ {
		<-c
	}
}

//通过waitgroup来让print输出完
func waitbyWG() {
	fmt.Println("waitbyWG")
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}
