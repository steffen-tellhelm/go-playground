package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const batchSize = 10
const waitSec = 1

type Message string

func writeBatch(msgBatch []Message, reason string) {
	log.Println(reason)
	log.Println(msgBatch)
}

func writer(ctx context.Context, msgCh <-chan Message) {
	timer := time.NewTicker(waitSec * time.Second)
	var msgBatch []Message

	go func() {
		for {
			select {
			case <-ctx.Done():
				writeBatch(msgBatch, "FINISH")
			case msg := <-msgCh:
				msgBatch = append(msgBatch, msg)
				if len(msgBatch) >= batchSize {
					writeBatch(msgBatch, "BATCH FULL")
					msgBatch = nil
				}
			case <-timer.C:
				if len(msgBatch) > 0 {
					writeBatch(msgBatch, "TIMEOUT")
					msgBatch = nil
				}
				timer.Reset(waitSec * time.Second)
			}
		}
	}()
}

func consumer(ctx context.Context) <-chan Message {
	msgCh := make(chan Message, 1)
	msgId := 0

	go func() {
		for {
			timespan := rand.Intn(10000)
			amount := rand.Intn(100) + 1

			select {
			case <-ctx.Done():
				close(msgCh)
				return
			case <-time.After(time.Duration(timespan) * time.Millisecond):
				log.Printf("consumed batch size: %v\n", amount)
				for i := 0; i < amount; i++ {
					msgId++
					msgCh <- Message(fmt.Sprintf("#%v", msgId))
				}
			}
		}
	}()

	return msgCh
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	msgCh := consumer(ctx)
	writer(ctx, msgCh)

	<-sigCh
	log.Println("\ngetting termination signal")
	cancel()
	log.Println("program finished")
}
