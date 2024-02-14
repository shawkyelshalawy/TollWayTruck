package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	//The connection used by the client
	conn *amqp.Connection
	// channels used to process/send messages
	ch *amqp.Channel
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}

func NewRabbitMQClient(conn *amqp.Connection) (RabbitClient, error) {
	ch, err := conn.Channel()
	if err != nil {
		return RabbitClient{}, err
	}
	// Puts the Channel in confirm mode, which will allow waiting for ACK or NACK from the receiver
	if err := ch.Confirm(false); err != nil {
		return RabbitClient{}, err
	}

	return RabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc RabbitClient) Close() error {
	return rc.ch.Close()
}

func (rc RabbitClient) CreateQueue(queueName string, durable, autodelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autodelete, false, false, nil)
	return err
}

// will bind the current channel to the given exchange based on  given cfgs
func (rc RabbitClient) CreateBinding(name, binding, exchange string) error {
	// leaving nowait set to false will make the channel return an error  if its fails to bind
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc RabbitClient) ProduceData(ctx context.Context, exchange, routingKey string, options amqp.Publishing) error {

	confirmation, err := rc.ch.PublishWithDeferredConfirmWithContext(ctx,
		exchange,
		routingKey,
		true,

		false,
		options,
	)
	if err != nil {
		return err
	}
	// Blocks until ACK from Server is receieved
	log.Println(confirmation.Wait())
	return nil
}

func (rc RabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}
