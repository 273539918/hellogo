package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	HttpHandler()
}

func HttpHandler() {
	ctx, cancel := NewContextWithTimeout()
	defer cancel()
	deal(ctx, cancel)
}

func NewContextWithTimeout() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*10)
}

func deal(ctx context.Context, cancel context.CancelFunc) {

	for i := 0; ; {
		time.Sleep(time.Second)
		fmt.Println("select...")
		select {
		case <-ctx.Done():
			fmt.Println("http request timeout")
			return
		default:
			fmt.Println("deal request ", i)
			i = i + 1
			//设置手动取消,不手动取消会在超时时间之后自动取消
			//if i > 3 {
			//	fmt.Println("trigger ctx cancel()")
			//	cancel()
			//}

		}
	}
}
