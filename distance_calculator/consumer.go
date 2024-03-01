package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shawkyelshalawy/TollWayTruck/aggregator/client"
	"github.com/shawkyelshalawy/TollWayTruck/types"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type RabbitConsumer struct {
	//The connection used by the client
	conn *amqp.Connection
	// channels used to process/send messages
	ch          *amqp.Channel
	calcService CalculatorServicer
	aggClient   client.Client
}

func ConnectRabbitMQ(username, password, host, vhost string) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
}
func NewRabbitConsumer(svc CalculatorServicer, aggClient client.Client) (RabbitConsumer, error) {
	conn, err := ConnectRabbitMQ("shawky", "secret", "localhost:5672", "tollway")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return RabbitConsumer{}, err
	}
	// Puts the Channel in confirm mode, which will allow waiting for ACK or NACK from the receiver
	if err := ch.Confirm(false); err != nil {
		return RabbitConsumer{}, err
	}
	return RabbitConsumer{
		conn:        conn,
		ch:          ch,
		calcService: svc,
		aggClient:   aggClient,
	}, nil
}

func (rc *RabbitConsumer) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}

func (rc *RabbitConsumer) ReadMessageLoop() {
	messageBus, err := rc.Consume("tollway-created", "distance-calculator", false)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	var blocking chan struct{}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(10)
	var data types.OBUData
	go func() {
		for message := range messageBus {
			msg := message
			g.Go(func() error {

				if err := json.Unmarshal(msg.Body, &data); err != nil {
					logrus.Errorf("Failed to unmarshal message")
					return err
				}
				distance, err := rc.calcService.CalculateDistance(data)
				if err != nil {
					logrus.Errorf("Failed to calculate distance")
					return err
				}
				req := &types.AggregateRequest{
					ObuID: int32(data.OBUID),
					Value: distance,
					Unix:  time.Now().UnixNano(),
				}
				if err := rc.aggClient.Aggregate(context.Background(), req); err != nil {
					logrus.Errorf("Failed to aggregate distance, error: %v", err)
				}
				logrus.Infof("Calculated distance: %.2f", distance)
				time.Sleep(10 * time.Second)
				return nil
			})
		}
	}()
	log.Println("consuming messages started, use CTRL+C to stop it")
	<-blocking
}

// func (rc *RabbitConsumer) ReadMessageLoop() {
// 	for rc.isRunning {
// 		messageBus, err := rc.Consume("tollway-created", "distance-calculator", false)
// 		if err != nil {
// 			panic(err)
// 		}
// 		var data types.OBUData
// 		//	for message := range messageBus {
// 		msg := <-messageBus
// 		if err := json.Unmarshal(msg.Body, &data); err != nil {
// 			logrus.Errorf("Failed to unmarshal message: %v", err)
// 			continue
// 		}
// 		distance, err := rc.calcService.CalculateDistance(data)
// 		if err != nil {
// 			logrus.Errorf("Failed to calculate distance: %v", err)
// 			continue
// 		}
// 		logrus.Infof("Calculated distance: %.2f", distance)
// 		req := types.Distance{
// 			OBUID: data.OBUID,
// 			Value: distance,
// 			Unix:  time.Now().UnixNano(),
// 		}
// 		if err := rc.aggClient.Aggregate(req); err != nil {
// 			logrus.Errorf("Failed to aggregate distance, error: %v", err)
// 			continue
// 		}

//			//	}
//		}
//		log.Println("consuming messages started, use CTRL+C to stop it")
//	}
func (rc *RabbitConsumer) Start() {
	logrus.Info("Starting RabbitMQ consumer")
	rc.ReadMessageLoop()
}
