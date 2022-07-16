package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	defer fmt.Println("a")
	defer fmt.Println("b")
	defer fmt.Println("c")

	loopfunc()

	fmt.Println("sleep start")
	time.Sleep(time.Second)
	fmt.Println("sleep end")

	//loopfunc()
}

func loopfunc() {
	lock := sync.Mutex{}
	for i := 0; i < 3; i++ {
		//go func(i int) {
		func(i int) {
			lock.Lock()
			//无论业务是否出错，defer都会执行
			defer lock.Unlock()
			fmt.Printf("loopfunc: %d \n", i)
		}(i)
	}
}
