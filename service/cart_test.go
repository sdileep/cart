package service

import (
	"testing"

	"github.com/sdileep/cart/service/entity"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

const (
	doveSoap = "Dove Soap"
)

func Test_cartService_AddProduct(t *testing.T) {
	type fields struct {
		catalog ProductService
	}
	type args struct {
		cartID    string
		productID string
		quantity  uint8
	}

	catalog := map[string]*entity.Product{
		doveSoap: {ID: doveSoap, Name: doveSoap, Price: 39.99},
	}
	defaultProductService := NewProductService(catalog)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Cart
		wantErr error
	}{
		{
			name: "user adds empty product",
			fields: fields{
				catalog: defaultProductService,
			},
			args: args{
				quantity:  5,
				productID: "",
				cartID:    "",
			},
			wantErr: preConditionError("productID", "empty"),
		},
		{
			name: "user adds unknown product",
			fields: fields{
				catalog: defaultProductService,
			},
			args: args{
				quantity:  5,
				productID: "unknown",
				cartID:    "",
			},
			wantErr: ProductNotFound,
		},
		{
			name: "user adds 'Dove Soap' product of '0' quantity",
			fields: fields{
				catalog: defaultProductService,
			},
			args: args{
				quantity:  0,
				productID: doveSoap,
				cartID:    "",
			},
			wantErr: preConditionError("quantity", "empty"),
		},
		{
			name: "user adds 'Dove Soap' product of '5' quantities",
			fields: fields{
				catalog: defaultProductService,
			},
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
			name: "user adds 'Dove Soap' product of '5' quantities to a nil cartID",
			fields: fields{
				catalog: defaultProductService,
			},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCartService(tt.fields.catalog)
			got, err := c.AddProduct(tt.args.cartID, tt.args.productID, tt.args.quantity)

			if tt.wantErr != nil {
				require.NotNil(t, err, "error expected")
				assert.Equal(t, tt.wantErr.Error(), err.Error(), "error")
				return
			}

			require.NoError(t, err, "error")

			require.EqualValues(t, tt.want.Items, got.Items, "items")
			assert.Equal(t, tt.want.Total, got.Total, "total")

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
			if got := c.computeTotal(tt.args.cart); got != tt.want {
				t.Errorf("cartService.computeTotal() = %v, want %v", got, tt.want)
			}
		})
	}
}
