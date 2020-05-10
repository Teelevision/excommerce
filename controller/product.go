package controller

import (
	"context"
	"errors"

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

// Get returns the requested product. ErrNotFound is returned if there is no
// product with the given id.
func (c *Product) Get(ctx context.Context, productID string) (*model.Product, error) {
	product, err := c.ProductRepository.FindProduct(ctx, productID)
	switch {
	case errors.Is(err, persistence.ErrNotFound):
		return nil, ErrNotFound
	case err == nil:
		return &model.Product{
			ID:    product.ID,
			Name:  product.Name,
			Price: product.Price,
		}, nil
	default:
		panic(err)
	}
}
