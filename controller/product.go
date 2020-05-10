package controller

import (
	"context"

	"github.com/Teelevision/excommerce/model"
	"github.com/Teelevision/excommerce/persistence"
)

// Product is the controller that handles products.
type Product struct {
	ProductRepository persistence.ProductRepository
}

// GetAll gets all products.
func (c *Product) GetAll(ctx context.Context) ([]*model.Product, error) {
	products, err := c.ProductRepository.FindAllProducts(ctx)
	switch {
	case err == nil:
		return products, nil
	default:
		panic(err)
	}
}
