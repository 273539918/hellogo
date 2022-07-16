package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type PCSequence struct {
	i     int
	mutex sync.Mutex
}

func (pcs *PCSequence) SafeAdd() int {
	pcs.mutex.Lock()
	defer pcs.mutex.Unlock()
	pcs.i = pcs.i + 1
	return pcs.i
}

//func (pcs *PCSequence) SafeRead() int {
//	pcs.mutex.Lock()
//	defer pcs.mutex.Unlock()
//	return pcs.i
//}

func main() {

	ch := make(chan int, 10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	pcsequence := &PCSequence{i: 0, mutex: sync.Mutex{}}
	defer cancel()
	defer close(ch)
	for i := 0; i < 2; i++ {
		go productor2(i, ch, ctx, pcsequence)
		go consumer2(i, ch)
	}

	select {
	case <-ctx.Done():
		fmt.Println("main over")
		//等待其他go也结束
		time.Sleep(time.Second * 5)
		return
	}

}

func productor2(id int, ch chan<- int, ctx context.Context, sequence *PCSequence) {
	for range time.Tick(time.Second) {
		select {
		case <-ctx.Done():
			fmt.Printf("productor %d is over\n", id)
			return
		default:
			value := sequence.SafeAdd()
			fmt.Printf("product %d to chan %d \n", id, value)
			ch <- value
		}
	}
}

func consumer2(id int, ch <-chan int) {
	for range time.Tick(time.Second) {
		if value, ok := <-ch; ok == true {
			fmt.Printf("consumer %d from chan %d \n", id, value)
		} else {
			fmt.Printf("consumer %d is over \n", id)
			return
		}
	}
}
