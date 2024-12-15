package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vladomdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/vladomdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/vladomdev/learn-pub-sub-starter/internal/routing"
)

func main() {
	fmt.Println("Starting Peril client...")

	// Connect to RabbitMQ Message broker
	rbtURL := "amqp://guest:guest@localhost:5672/"
	rbtConnect, err := amqp.Dial(rbtURL)
	if err != nil {
		log.Fatalf("An error occured while connecting to RabbitMQ: %v", err)
	}
	fmt.Println("Successfully connected to RabbitMQ server.")
	defer rbtConnect.Close()

	// Log in to the user
	loggedUser := ""
	for {
		userName, err := gamelogic.ClientWelcome()
		if err == nil {
			loggedUser = userName
			break
		}
	}

	// Declare and bind a transient queue
	queueName := routing.PauseKey + "." + loggedUser
	_, _, err = pubsub.DeclareAndBind(rbtConnect, routing.ExchangePerilDirect, queueName, routing.PauseKey, 1)
	if err != nil {
		log.Fatalf("An error occured while Declaring and binding a Qeue: %v", err)
	}

	// wait for ctrl+c
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	signal.Stop(sigChan)
	close(sigChan)
	fmt.Println("\nDisconnecting from Peril server...")
}
