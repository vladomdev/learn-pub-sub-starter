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

	gameState := gamelogic.NewGameState(loggedUser)

	for {
		userInput := gamelogic.GetInput()
		if len(userInput) == 0 {
			continue
		}
		if userInput[0] == "spawn" {
			err = gameState.CommandSpawn(userInput)
			if err != nil {
				fmt.Printf("An error occured while spawning %v's units: %v", loggedUser, err)
			}
			continue
		}
		if userInput[0] == "move" {
			_, err = gameState.CommandMove(userInput)
			if err != nil {
				fmt.Printf("An error occured while moving %v's units: %v", loggedUser, err)
			}
			continue
		}
		if userInput[0] == "status" {
			gameState.CommandStatus()
			continue
		}
		if userInput[0] == "help" {
			gamelogic.PrintClientHelp()
			continue
		}
		if userInput[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
			continue
		}
		if userInput[0] == "quit" {
			gamelogic.PrintQuit()
			break
		}
		fmt.Println("\nCommand not understood. Please enter a valid command. Use 'help' to list valid commands.")
		continue
	}

	/*
		// wait for ctrl+c
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		signal.Stop(sigChan)
		close(sigChan)
		fmt.Println("\nDisconnecting from Peril server...")
	*/
}
