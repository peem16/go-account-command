package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"go-account-command/handler"
	"go-account-command/messageBroker"
	"go-account-command/repository"
	"go-account-command/router"
	"go-account-command/services"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading environment: %v", err)
	}
}

func initDatabase() *mongo.Database {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("DSN")))
	if err != nil {
		panic("failed to connect database")
	}
	database := client.Database("account")

	return database
}

func initKafka() (sarama.SyncProducer, sarama.ConsumerGroup) {
	brokers := []string{os.Getenv("BROKER")}

	config := sarama.NewConfig()
	config.Net.SASL.User = os.Getenv("KAFKA_USER")
	config.Net.SASL.Password = os.Getenv("KAFKA_PASSWORD")
	config.Net.SASL.Handshake = true
	config.Net.SASL.Enable = true
	config.Net.TLS.Enable = true
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ClientAuth:         0,
	}
	config.Net.TLS.Config = tlsConfig
	config.Producer.Timeout = 5 * time.Second
	config.Producer.Return.Successes = true

	client, errClient := sarama.NewClient(brokers, config)
	if errClient != nil {
		log.Printf("sarama.NewClient err, message=%s \n", errClient)
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		log.Printf("sarama.NewSyncProducerFromClient err, message=%s \n", err)
	}
	consumer, err := sarama.NewConsumerGroupFromClient("Group-GO", client)
	if err != nil {
		log.Printf("sarama.NewConsumerGroupFromClient err, message=%s \n", err)
	}

	return producer, consumer
}

func main() {

	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")
	db := initDatabase()
	producer, consumer := initKafka()
	defer producer.Close()
	defer consumer.Close()

	eventProducer := messageBroker.NewEventProducer(producer)

	accountEventRepository := repository.NewAccountEventRepositoryDB(db)
	balanceViweRepository := repository.NewBalanceViweRepositoryDB(db)
	accountEventService := services.NewAccountEventService(accountEventRepository, balanceViweRepository, eventProducer)
	balanceViewService := services.NewBalanceViewService(balanceViweRepository, accountEventRepository)

	accountEventHandler := handler.NewAccountEventHandler(accountEventService)
	consumerHandler := messageBroker.NewConsumerHandler(balanceViewService)

	r := router.New()

	r.POST("api/v1/account", accountEventHandler.AccountEventCreate)
	r.POST("api/v1/account/transfer", accountEventHandler.AccountEventTransfer)
	r.DEL("api/v1/account/clear", accountEventHandler.ClearAccount)

	ctx, cancel := context.WithCancel(context.Background())

	go r.ListenAndServe(ctx, cancel)()
	messageBroker.ConsumerListen(consumer, consumerHandler, ctx, cancel)
	fmt.Println("start service")
}
