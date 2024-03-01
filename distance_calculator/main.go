package main

import "github.com/shawkyelshalawy/TollWayTruck/aggregator/client"

const (
	aggregatorEndpoint = "http://localhost:8080/aggregate"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)

	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	// grpcClient, err := client.NewGRPCClient(aggregatorEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	rcconsumer, err := NewRabbitConsumer(svc, httpClient)
	if err != nil {
		panic(err)
	}
	rcconsumer.ReadMessageLoop()
}
