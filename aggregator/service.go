package main

import (
	"fmt"

	"github.com/shawkyelshalawy/TollWayTruck/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}
func (a *InvoiceAggregator) AggregateDistance(data types.Distance) error {
	fmt.Println("Aggregating distance")
	return a.store.Insert(data)
}
