package service

import (
	"errors"

	"github.com/sdileep/cart/service/entity"
)

var (
	ProductNotFound = errors.New("product not found")
)

type ProductService interface {
	Lookup(productID string) (*entity.Product, error)
}

type productService struct {
	catalog map[string]*entity.Product
}

func (p *productService) Lookup(productID string) (*entity.Product, error) {
	if len(p.catalog) == 0 {
		return nil, ProductNotFound
	}

	product, ok := p.catalog[productID]
	if !ok {
		return nil, ProductNotFound
	}

	return product, nil
}

func NewProductService(catalog map[string]*entity.Product) ProductService {
	if catalog == nil {
		catalog = make(map[string]*entity.Product)
	}
	return &productService{catalog}
}
