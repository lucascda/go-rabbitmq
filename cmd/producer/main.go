package main

import (
	"log"
	"os"
	"time"

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

	time.Sleep(10 * time.Second)
	log.Println(client)
}
