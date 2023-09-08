package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// TestContext 测试
func TestContext() {
	root := context.Background()
	ctx, fn := context.WithTimeout(root, 10*time.Second)
	ctx2, cancel := context.WithDeadline(root, time.Now().Add(4*time.Second))
	defer cancel()
	defer fn()
	group := sync.WaitGroup{}
	group.Add(1)
	ch := make(chan int32, 2)
	go func(ctx context.Context) {
		defer group.Done()
		fmt.Println("start")
		i := 0
		for {
			select {
			case <-time.After(1 * time.Second):
				fmt.Println(i)
				i++
			case <-ctx.Done():
				fmt.Println("canceled")
				ch <- 1
				return
			}
		}
	}(ctx2)
	select {
	case <-ch:
		fmt.Println("top")
	case <-ctx.Done():
		fmt.Println("boom!")
	}
	group.Wait()
}

func TestTimer() {

}

func main() {
	TestContext()
}
