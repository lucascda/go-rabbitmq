package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/lucas_cda/go-rabbitmq/internal"
	"golang.org/x/sync/errgroup"
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

	ctx := context.Background()

	var blocking chan struct{}

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)

	go func() {

		for message := range messageBus {
			msg := message
			g.Go(func() error {
				log.Printf("New message: %v", msg)
				time.Sleep(10 * time.Second)
				if err := msg.Ack(false); err != nil {
					log.Println("ack msg failed")
					return err
				}
				log.Printf("Acknowledged msg %s", message.MessageId)
				return nil

			})
		}
	}()
	log.Println("consuming...")
	<-blocking

}
