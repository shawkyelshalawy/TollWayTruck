package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/shawkyelshalawy/TollWayTruck/types"
)

func main() {
	listenAddr := flag.String("listen-addr", ":8080", "server listen address")
	flag.Parse()
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewLogMiddleware(svc)
	makeHttpTransport(*listenAddr, svc)
}

func makeHttpTransport(listenAddr string, svc Aggregator) {
	fmt.Println("Starting HTTP transport on", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)

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
