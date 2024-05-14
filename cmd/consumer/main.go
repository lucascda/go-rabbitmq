package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/lucas_cda/go-rabbitmq/internal"
)

func main() {
	godotenv.Load()
	rabbitmq := struct {
		user  string
		pass  string
		host  string
		port  string
		vhost string
	}{
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASS"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"),
		os.Getenv("RABBITMQ_VHOST"),
	}
	conn, err := internal.ConnectRabbitMQ(rabbitmq.user, rabbitmq.pass, rabbitmq.host, rabbitmq.port, rabbitmq.vhost)

	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client, err := internal.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	messageBus, err := client.Consume("customers_created", "email-service", false)
	if err != nil {
		panic(err)
	}

	var blocking chan struct{}

	go func() {
		for message := range messageBus {
			log.Println("New message: %v", message)

			if err := message.Ack(false); err != nil {
				log.Println("acknowledge message fail")
				continue
			}
			log.Printf("Acknowledge message %s\n", message.MessageId)
		}
	}()

	fmt.Println("Consuming...")
	<-blocking
}
