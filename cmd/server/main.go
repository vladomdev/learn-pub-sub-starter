package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/vladomdev/learn-pub-sub-starter/internal/gamelogic"
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

	// Print available commands for users
	gamelogic.PrintServerHelp()

	// Get user input
	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		if userInput[0] == "pause" {
			fmt.Println("\nPausing the game.")
			playingState.IsPaused = true
			err = pubsub.PublishJSON(rbtChannel, routing.ExchangePerilDirect, routing.PauseKey, playingState)
			if err != nil {
				log.Fatalf("Failed to publish message to Exchange: %v", err)
			}
			continue
		}
		if userInput[0] == "resume" {
			fmt.Println("\nResuming the game.")
			playingState.IsPaused = false
			err = pubsub.PublishJSON(rbtChannel, routing.ExchangePerilDirect, routing.PauseKey, playingState)
			if err != nil {
				log.Fatalf("Failed to publish message to Exchange: %v", err)
			}
			continue
		}
		if userInput[0] == "quit" {
			fmt.Println("Shutting down Peril server...")
			break
		}
		if userInput[0] == "help" {
			gamelogic.PrintServerHelp()
			continue
		}
		fmt.Println("\nCommand not understood. Please enter a valid command. Use 'help' to list valid commands.")
		continue
	}
}
