package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func opserver(cancel func(), pingCh chan int) {
	pingCh <- 1
	time.Sleep(10 * time.Second)
	cancel()
}

func ping(ctx context.Context, pingCh chan int, pongCh chan float32) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\nPing: I'm done\n")
			close(pingCh)
			return
		case val := <-pongCh:
			time.Sleep(1 * time.Second)
			fmt.Println(" --> Ping:", val)
			pingVal := rand.Intn(10)
			fmt.Print("Ping: ", pingVal)
			time.Sleep(200 * time.Millisecond)
			pingCh <- pingVal
		}
	}
}

func pong(ctx context.Context, pingCh chan int, pongCh chan float32) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\nPong: I'm done\n")
			close(pongCh)
			return
		case val := <-pingCh:
			fmt.Println(" --> Pong:", val)
			pongVal := rand.Float32()
			fmt.Print("Pong: ", pongVal)
			time.Sleep(200 * time.Millisecond)
			pongCh <- pongVal
		}
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	pingCh := make(chan int)
	pongCh := make(chan float32)

	go ping(ctx, pingCh, pongCh)
	go pong(ctx, pingCh, pongCh)
	go opserver(cancel, pingCh)

	<-ctx.Done()
	time.Sleep(2 * time.Second)
}
