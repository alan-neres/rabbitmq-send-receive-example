package main

import (
	"fmt"
	"os"

	"github.com/streadway/amqp" // Importação da biblioteca AMQP
)

func main() {
	rmqURL := os.Getenv("RMQ_URL")
	if rmqURL == "" {
		fmt.Println("Variável de ambiente RMQ_URL não definida")
		return
	}

	// Conexão com o RabbitMQ
	conn, err := amqp.Dial(rmqURL)
	if err != nil {
		fmt.Println("Falha ao conectar com o broker", err)
		return
	}
	defer conn.Close()

	// Criando um canal
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Falha ao abrir um canal com o broker", err)
		return
	}
	defer ch.Close()

	// Criação de uma Queue // Caso já exista, simplesmente se conecta
	q, err := ch.QueueDeclare(
		"cobrar", // Nome da fila
		true,     // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		fmt.Println("Falha ao criar a queue", err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			fmt.Printf("Mensagem de cobrança recebida na queue %v: %v\n", q.Name, string(d.Body))
			d.Ack(true)
		}
	}()

	fmt.Println("[Cobranca de Vendas] Aguardando por mensagens")
	<-forever
}
