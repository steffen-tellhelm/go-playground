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

func writer(ctx context.Context, msgCh <-chan []Message) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("finish writing")
			case msgBatch := <-msgCh:
				log.Printf("writing batch (size: %v): %v\n", len(msgBatch), msgBatch)
			}
		}
	}()
}

func consumer(ctx context.Context) <-chan []Message {
	msgCh := make(chan []Message, 1)
	var msgBatch []Message
	msgId := 0
	timer := time.NewTicker(waitSec * time.Second)

	go func() {
		for {
			timespan := rand.Intn(10000)
			amount := rand.Intn(100) + 1

			select {
			case <-ctx.Done():
				close(msgCh)
				return
			case <-timer.C:
				if len(msgBatch) > 0 {
					log.Println("write reason: TIMEOUT")
					msgCh <- msgBatch
					msgBatch = nil
				}
				timer.Reset(waitSec * time.Second)
			case <-time.After(time.Duration(timespan) * time.Millisecond):
				log.Printf("consumed batch size: %v\n", amount)
				for i := 0; i < amount; i++ {
					msgId++
					msg := Message(fmt.Sprintf("#%v", msgId))
					msgBatch = append(msgBatch, msg)
					if len(msgBatch) >= batchSize {
						log.Println("write reason: BATCH FULL")
						msgCh <- msgBatch
						msgBatch = nil
					}
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
