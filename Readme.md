# Tollwaytruck

- A microservices app that simulates a tollway system for trucks calculates the distance that the tollway truck has passed and calculates the tolls for it, so you can get the invoice for the truck before it does the trip.

## Services

`obu` : a service that simulates the on board unit of the truck , it sends data ( longitude & latitude) to the datareceiver service.

`datarecevier` : a service that receives data from the trucks and sends it to the data distance_calculator service.

`distance_calculator` : a service that calculates the distance between two points.

`aggregator`: a service that aggregates the data from the distance_calculator , calulates the tolls for the trucks and sends the data to the data `storage service (todo)`.

`gateway`: a service that acts as a gateway for the services , it communicates with the aggregator service , sends the id of the truck and get the invoice for it

## Technologies used in this project:

- Golang
- Rabbitmq
- GRPC
- Prometheus
- Docker
- Grafana

## Rabbitmq

**get rabbitmq container**

``
docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management

``

**Add user and password and make him administrator**

``

- docker exec rabbitmq rabbitmqctl add_user shawky secret
- docker exec rabbitmq rabbitmqctl set_user_tags shawky administrator
  ``

**Add Vhost**

``
docker exec rabbitmq rabbitmqctl add_vhost tollway

``

## Installing GRPC and Protobuffer plugins for Golang.

1. Protobuffers

`go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28`

2. GRPC

`go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2`

## Installing Prometheus

Install Prometheus in a Docker container

`docker run -p 9090:9090 -v ./.config/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus`

Installing prometheus golang client

`go get github.com/prometheus/client_golang/prometheus`

## Installing Grafana

`docker run -d -p 3000:3000 grafana/grafana`
