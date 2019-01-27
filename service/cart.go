package service

import (
	"fmt"
	"math"
	"time"

	"github.com/sdileep/cart/service/entity"
)

func preConditionError(attribute, message string) error {
	return fmt.Errorf("pre-condition failed, %s: %s", attribute, message)
}

type CartService interface {
	AddProduct(cartID string, productID string, quantity uint8) (*entity.Cart, error)
}

type cartService struct {
	// mocking cart data store with activeCarts for simplicity of the exercise
	activeCarts    map[string]*entity.Cart
	productService ProductService
	taxService     TaxService
}

func (c *cartService) AddProduct(cartID string, productID string, quantity uint8) (*entity.Cart, error) {
	if productID == "" {
		return nil, preConditionError("productID", "empty")
	}

	if quantity == 0 {
		return nil, preConditionError("quantity", "empty")
	}

	// lookup product
	product, err := c.productService.Lookup(productID)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ProductNotFound
	}

	cart, ok := c.activeCarts[cartID]
	if cart == nil || !ok {
		// keeping ID generation simple intentionally
		cartID := string(time.Now().Unix())
		cart = &entity.Cart{ID: cartID}
		c.activeCarts[cartID] = cart
	}

	var item *entity.CartItem
	for _, v := range cart.Items {
		if productID == v.ProductID {
			item = v
			item.Quantity += quantity
			break
		}
	}

	if item == nil {
		item = &entity.CartItem{
			ProductID: product.ID,
			Quantity:  quantity,
			UnitPrice: product.Price,
		}
		cart.Items = append(cart.Items, item)
	}

	total := c.computeTotal(cart)
	tax := c.computeTax(total)

	cart.Tax = math.Ceil(tax*100) / 100
	cart.Total = math.Ceil((total+tax)*100) / 100

	return cart, nil
}

func (c *cartService) computeTotal(cart *entity.Cart) float64 {
	if cart == nil {
		return 0
	}

	// loop through cart items to calculate total
	var total float64
	for _, v := range cart.Items {
		total += float64(v.Quantity) * v.UnitPrice
	}

	return total
}

func (c *cartService) computeTax(total float64) float64 {
	return c.taxService.ComputeTax(total)
}

func NewCartService(productService ProductService, taxService TaxService) CartService {
	return &cartService{
		activeCarts:    make(map[string]*entity.Cart),
		productService: productService,
		taxService:     taxService,
	}
}
