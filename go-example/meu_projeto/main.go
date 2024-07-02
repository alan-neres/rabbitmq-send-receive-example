package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Configuração da conexão RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Abertura de um canal
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declaração de uma fila para enviar mensagens
	q, err := ch.QueueDeclare(
		"hello", // nome da fila
		false,   // durável
		false,   // exclui quando não usado
		false,   // exclusiva
		false,   // sem argumentos extras
		nil,     // argumentos extras
	)
	failOnError(err, "Failed to declare a queue")

	// Mensagem a ser enviada
	body := "Hello, World!"

	// Usando time.Tick para enviar mensagem a cada segundo
	ticker := time.Tick(500 * time.Millisecond)

	for range ticker {
		// Publicar a mensagem na fila
		err := ch.Publish(
			"hello", // exchange
			q.Name,  // routing key (nome da fila)
			false,   // mandatory
			false,   // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")

		fmt.Println(" [x] Sent 'Hello World!'")
	}
}
