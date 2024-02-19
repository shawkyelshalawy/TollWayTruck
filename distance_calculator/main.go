package main

import "github.com/shawkyelshalawy/TollWayTruck/aggregator/client"

const (
	aggregatorEndpoint = "http://localhost:8080/aggregate"
)

func main() {
	var (
		svc CalculatorServicer
		err error
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	httpClient := client.NewClient(aggregatorEndpoint)
	rcconsumer, err := NewRabbitConsumer(svc, httpClient)
	if err != nil {
		panic(err)
	}
	rcconsumer.ReadMessageLoop()
}
