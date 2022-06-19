package messageBroker

import (
	"context"
	"go-account-command/events"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
)

//sarama-kafka
func ConsumerListen(consumer sarama.ConsumerGroup, consumerHandler sarama.ConsumerGroupHandler, ctx context.Context, cancel context.CancelFunc) {
	keepRunning := true
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {

			if err := consumer.Consume(ctx, events.Topics, consumerHandler); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	for keepRunning {
		select {
		case <-ctx.Done():
			log.Println("terminating: context cancelled")
			keepRunning = false
		case <-sigterm:
			log.Println("terminating: via signal")
			keepRunning = false
		}
	}
	cancel()
	wg.Wait()
	if err := consumer.Close(); err != nil {
		log.Panicf("Error closing consumer: %v", err)
	}

}
