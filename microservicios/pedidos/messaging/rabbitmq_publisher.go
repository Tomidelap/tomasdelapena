package messaging

import (
	"encoding/json"

	amqp "github.com/rabbitmq/amqp091-go"
)

const ColaPedidosConfirmados = "pedidos-confirmados"

type RabbitMQPublisher struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMQPublisher(url string) (*RabbitMQPublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Crea la cola si todavía no existe (durable).
	if _, err := ch.QueueDeclare(ColaPedidosConfirmados, true, false, false, false, nil); err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}

	return &RabbitMQPublisher{conn: conn, ch: ch}, nil
}

func (p *RabbitMQPublisher) Publicar(evento any) error {
	body, err := json.Marshal(evento)
	if err != nil {
		return err
	}
	return p.ch.Publish("", ColaPedidosConfirmados, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
}

func (p *RabbitMQPublisher) Close() {
	p.ch.Close()
	p.conn.Close()
}
