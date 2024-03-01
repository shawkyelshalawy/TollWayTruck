package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/shawkyelshalawy/TollWayTruck/types"
)

type APIError struct {
	Code int
	Err  error
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
		//grpcListenAddr = os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAddr = os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewLogMiddleware(svc)
	makeHttpTransport(httpListenAddr, svc)
}

// func makeGRPCTransport(listenAddr string, svc Aggregator) error {
// 	fmt.Println("GRPC transport running on port ", listenAddr)
// 	// Make a TCP listener
// 	ln, err := net.Listen("tcp", listenAddr)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		fmt.Println("stopping GRPC transport")
// 		ln.Close()
// 	}()
// 	// Make a new GRPC native server with (options)
// 	server := grpc.NewServer([]grpc.ServerOption{}...)
// 	// Register GRPC server implementation to the GRPC package.
// 	types.RegisterAggregatorServer(server, NewAggregatorGRPCServer(svc))
// 	return server.Serve(ln)
// }

func makeHttpTransport(listenAddr string, svc Aggregator) {
	fmt.Println("Starting HTTP transport on", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice(svc))
	http.ListenAndServe(listenAddr, nil)

}

func handleInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing obuId "})
			return
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obuID"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuId)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data types.Distance
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		}
		if err := svc.AggregateDistance(data); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
