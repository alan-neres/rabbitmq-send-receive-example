package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
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

	// Criação da Exchange
	err = ch.ExchangeDeclare(
		"ecommerce.nova.venda", // Nome da exchange
		"direct",               // Tipo da exchange
		true,                   // durable
		false,                  // auto-deleted
		false,                  // internal
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		fmt.Println("Falha ao construir a exchange", err)
		return
	}

	// Criação de uma Queue
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

	// Associando a Queue até a Exchange
	// e informando a binding key para roteamento
	err = ch.QueueBind(
		q.Name,                 // Nome da fila
		"cobrar",               // Binding key de roteamento
		"ecommerce.nova.venda", // Nome da exchange
		false,                  // no-wait
		nil,                    // arguments
	)
	if err != nil {
		fmt.Println("Falha ao realizar o binding da queue com a exchange", err)
		return
	}

	fmt.Println("Configuração concluída com sucesso")

	// Usando time.Tick para enviar mensagens periodicamente
	ticker := time.Tick(5 * time.Second)

	for i := 0; i < 100; i++ {
		id := uuid.New()

		// Mensagem simples
		body := fmt.Sprintf("id:%v", id)

		// Publicando a mensagem na exchange
		err = ch.Publish(
			"ecommerce.nova.venda", // exchange
			"cobrar",               // routing key (binding key)
			false,                  // mandatory
			false,                  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			fmt.Println("Falha ao publicar a mensagem na exchange", err)
		}

		fmt.Printf("Mensagem de venda enviada para a exchange ecommerce.nova.venda: %v\n", body)

		// Espera pelo próximo tick do ticker
		<-ticker
	}
}
