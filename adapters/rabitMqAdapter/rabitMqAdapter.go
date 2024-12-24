package rabitMqAdapter

import (
	"cdc-file-processor/domain/ports"
	"context"
    "log"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMqAdapter struct {
	conn *amqp.Connection
	ch *amqp.Channel
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NewRabbitMqAdapter() ports.QueueSp {
	return &RabbitMqAdapter{
		conn: nil,
		ch: nil,
	}
}

func (r *RabbitMqAdapter) SendMessage(message string, qName string) error {
	if r.conn == nil || r.conn.IsClosed() {
        err := r.Connect()
        if err != nil {
            return err 
        }
    }

	// Declare a queue that will be created if not exists with some args
	q, err := r.ch.QueueDeclare(
		qName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := message
	err = r.ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
	return nil
}

func (r *RabbitMqAdapter) ReceiveMessage() (string, error) {
	return "", nil
}

func (r *RabbitMqAdapter) Connect() error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	r.conn = conn

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	r.ch = ch

	return nil
}

func (r *RabbitMqAdapter) Close() error {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.ch != nil {
		r.ch.Close()
	}
	return nil
}