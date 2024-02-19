package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/shawkyelshalawy/TollWayTruck/rabbitmq"
	"github.com/shawkyelshalawy/TollWayTruck/types"
)

type DataReceiver struct {
	msgch        chan types.OBUData
	conn         *websocket.Conn
	rabbitClient rabbitmq.RabbitClient
}

func main() {
	conn, err := rabbitmq.ConnectRabbitMQ("shawky", "secret", "localhost:5672", "tollway")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client, err := rabbitmq.NewRabbitMQClient(conn)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	if err := client.CreateQueue("tollway-created", true, false); err != nil {
		panic(err)
	}
	if err := client.CreateBinding("tollway-created", "tollway.*", "tollway_events"); err != nil {
		panic(err)
	}
	recv := NewDataReceiver(client)
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}
func NewDataReceiver(client rabbitmq.RabbitClient) *DataReceiver {
	return &DataReceiver{
		msgch:        make(chan types.OBUData, 128),
		rabbitClient: client,
	}
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU connected client connected !")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		//fmt.Printf("received OBU data from[%d]:: <lat %.2f, long %.2f> \n", data.OBUID, data.Lat, data.Long)
		dataBytes, err := json.Marshal(data)
		if err != nil {
			log.Println("Failed to marshal OBU data:", err)
			continue
		}
		pub := amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         dataBytes,
		}

		if err := dr.rabbitClient.ProduceData(ctx, "tollway_events", "tollway.egy", pub); err != nil {
			log.Println("Failed to send message:", err)
			continue
		}
		fmt.Println("received message")
	}
}
