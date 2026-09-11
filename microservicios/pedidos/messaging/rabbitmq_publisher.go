package messaging

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"pedidos/models"
)

const colaPedidosConfirmados = "pedidos-confirmados"

type Publisher interface {
	PublicarPedidoConfirmado(evento models.PedidoConfirmado) error
}

type RabbitMQPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NuevoRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	_, err = channel.QueueDeclare(colaPedidosConfirmados, true, false, false, false, nil)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublisher{conn: conn, channel: channel}, nil
}

func (p *RabbitMQPublisher) PublicarPedidoConfirmado(evento models.PedidoConfirmado) error {
	cuerpo, err := json.Marshal(evento)
	if err != nil {
		return err
	}

	err = p.channel.Publish("", colaPedidosConfirmados, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        cuerpo,
	})
	if err != nil {
		return err
	}

	log.Printf("evento publicado en %s: %s", colaPedidosConfirmados, cuerpo)
	return nil
}
