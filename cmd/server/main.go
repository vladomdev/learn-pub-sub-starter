package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vladomdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/vladomdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril server...")

	// Connect to RabbitMQ Message broker
	rbtURL := "amqp://guest:guest@localhost:5672/"
	rbtConnect, err := amqp.Dial(rbtURL)
	if err != nil {
		log.Fatalf("An error occured while connecting to RabbitMQ: %v", err)
	}
	fmt.Println("Successfully connected to RabbitMQ server.")
	defer rbtConnect.Close()

	// Create a new channel
	rbtChannel, err := rbtConnect.Channel()
	if err != nil {
		log.Fatalf("An error occured while creating a new RabbitMQ Channel: %v", err)
	}

	// Publish a message to the exchange
	playingState := routing.PlayingState{
		IsPaused: true,
	}
	err = pubsub.PublishJSON(rbtChannel, routing.ExchangePerilDirect, routing.PauseKey, playingState)
	if err != nil {
		log.Fatalf("Failed to publish message to Exchange: %v", err)
	}

	// wait for ctrl+c
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	signal.Stop(sigChan)
	close(sigChan)
	fmt.Println("\nShutting down Peril server...")
}
