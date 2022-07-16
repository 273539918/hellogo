package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go Speak(ctx, cancel)
	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(1 * time.Second)
}

func Speak(ctx context.Context, cancel context.CancelFunc) {
	for v := range time.Tick(time.Second) {
		fmt.Println(v)
		select {
		case <-ctx.Done():
			fmt.Println("over")
			return
		default:
			fmt.Println("speak something")
		}
	}
}
