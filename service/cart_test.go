package service

import (
	"math"
	"testing"

	"github.com/sdileep/cart/service/entity"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

const (
	axeDeo   = "Axe Deo"
	doveSoap = "Dove Soap"
)

func Test_cartService_AddProduct(t *testing.T) {
	const taxRate = 12.5
	type fields struct {
		productService ProductService
		taxService     TaxService
	}
	type args struct {
		cartID    string
		productID string
		quantity  uint8
	}

	productMaster := map[string]*entity.Product{
		axeDeo:   {ID: axeDeo, Name: axeDeo, Price: 99.99},
		doveSoap: {ID: doveSoap, Name: doveSoap, Price: 39.99},
	}
	defaultProductService := NewProductService(productMaster)
	defaultTaxService := NewTaxService(0)
	defaultFields := fields{
		productService: defaultProductService,
		taxService:     defaultTaxService,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		setup   func() []args
		want    *entity.Cart
		wantErr error
	}{
		{
			name:   "user adds empty product",
			fields: defaultFields,
			args: args{
				quantity:  5,
				productID: "",
				cartID:    "",
			},
			wantErr: preConditionError("productID", "empty"),
		},
		{
			name:   "user adds unknown product",
			fields: defaultFields,
			args: args{
				quantity:  5,
				productID: "unknown",
				cartID:    "",
			},
			wantErr: ProductNotFound,
		},
		{
			name:   "user adds 'Dove Soap' product of '0' quantity",
			fields: defaultFields,
			args: args{
				quantity:  0,
				productID: doveSoap,
				cartID:    "",
			},
			wantErr: preConditionError("quantity", "empty"),
		},
		{
			name:   "user adds 'Dove Soap' product of '5' quantities",
			fields: defaultFields,
			args: args{
				quantity:  5,
				productID: doveSoap,
				cartID:    "",
			},
			want: &entity.Cart{
				Items: []*entity.CartItem{
					{ProductID: doveSoap, Quantity: 5, UnitPrice: 39.99},
				},
				Total: 199.95,
			},
		},
		{
			name:   "user adds 'Dove Soap' product of '5' quantities to a nil cart",
			fields: defaultFields,
			args: args{
				quantity:  5,
				productID: doveSoap,
			},
			want: &entity.Cart{
				Items: []*entity.CartItem{
					{ProductID: doveSoap, Quantity: 5, UnitPrice: 39.99},
				},
				Total: 199.95,
			},
		},
		{
			name: `user adds 'Dove Soap' product of '5' quantities
					and  'Dove Soap' product of '3' quantities
					`,
			fields: defaultFields,
			args: args{
				quantity:  3,
				productID: doveSoap,
			},
			setup: func() []args {
				return []args{
					{
						quantity:  5,
						productID: doveSoap,
					},
				}

			},
			want: &entity.Cart{
				Items: []*entity.CartItem{
					{ProductID: doveSoap, Quantity: 8, UnitPrice: 39.99},
				},
				Total: 319.92,
			},
		},
		{
			name: "2 'Dove Soap', 2 'Axe Deo', tax rate '12.5%'",
			fields: fields{
				productService: defaultProductService,
				taxService:     NewTaxService(taxRate),
			},
			args: args{
				quantity:  2,
				productID: axeDeo,
			},
			setup: func() []args {
				return []args{
					{
						quantity:  2,
						productID: doveSoap,
					},
				}

			},
			want: &entity.Cart{
				Items: []*entity.CartItem{
					{ProductID: doveSoap, Quantity: 2, UnitPrice: 39.99},
					{ProductID: axeDeo, Quantity: 2, UnitPrice: 99.99},
				},
				Tax:   35.00,
				Total: 314.96,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCartService(tt.fields.productService, tt.fields.taxService)

			var cart *entity.Cart
			var cartID string
			var err error
			if tt.setup != nil {
				argsArr := tt.setup()
				for _, args := range argsArr {

					cart, err = c.AddProduct(cartID, args.productID, args.quantity)
					cartID = cart.ID
				}
			}

			got, err := c.AddProduct(cartID, tt.args.productID, tt.args.quantity)

			if tt.wantErr != nil {
				require.NotNil(t, err, "error expected")
				assert.Equal(t, tt.wantErr.Error(), err.Error(), "error")
				return
			}

			require.NoError(t, err, "error")

			require.EqualValues(t, tt.want.Items, got.Items, "items")
			assert.Equal(t, tt.want.Total, got.Total, "total")
			assert.Equal(t, tt.want.Tax, got.Tax, "tax")

		})
	}
}

func Test_cartService_computeTotal(t *testing.T) {
	type fields struct {
		catalog ProductService
	}
	type args struct {
		cart *entity.Cart
	}

	catalog := map[string]*entity.Product{
		doveSoap: {ID: doveSoap, Name: doveSoap, Price: 39.99},
	}
	defaultProductService := NewProductService(catalog)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "nil cart",
			fields: fields{
				catalog: defaultProductService,
			},
			args: args{},
			want: 0,
		},
		{
			name: "empty cart",
			fields: fields{
				catalog: defaultProductService,
			},
			args: args{
				cart: &entity.Cart{},
			},
			want: 0,
		},
		{
			name: "cart with items",
			fields: fields{
				catalog: defaultProductService,
			},
			args: args{
				cart: &entity.Cart{
					Items: []*entity.CartItem{
						{ProductID: doveSoap, Quantity: 5, UnitPrice: 39.99},
					},
				},
			},
			want: 199.95,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cartService{
				productService: tt.fields.catalog,
			}
			got := c.computeTotal(tt.args.cart)
			got = math.Ceil(got*100) / 100
			if got != tt.want {
				t.Errorf("cartService.computeTotal() = %v, want %v", got, tt.want)
			}
		})
	}
}
