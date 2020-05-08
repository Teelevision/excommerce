package controller

import (
	"context"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
)

// GetProduct is the controller that gets products.
type GetProduct struct {
	ProductRepository persistence.ProductRepository
}

// All gets all products.
func (c *GetProduct) All(ctx context.Context) ([]*model.Product, error) {
	products, err := c.ProductRepository.FindAllProducts(ctx)
	switch {
	case err == nil:
		return products, nil
	default:
		panic(err)
	}
}
