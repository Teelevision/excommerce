package log

import (
	"context"
	"log"

	"github.com/Teelevision/excommerce/persistence"
)

// Adapter is a write only persistence adapter. It simply prints everything to
// the standard output.
type Adapter struct{}

// NewAdapter returns a new log adapter.
func NewAdapter() *Adapter {
	return &Adapter{}
}

var _ persistence.PlacedOrderRepository = (*Adapter)(nil)

// PlaceOrder places the order and all related data.
func (a *Adapter) PlaceOrder(ctx context.Context, order persistence.PlacedOrder) error {
	log.Print(order)
	return nil
}
