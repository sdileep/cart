package service

import (
	"testing"

	"github.com/sdileep/cart/service/entity"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func Test_productCatalog_Lookup(t *testing.T) {
	type fields struct {
		master map[string]*entity.Product
	}
	type args struct {
		productID string
	}

	productMaster := map[string]*entity.Product{
		doveSoap: {ID: doveSoap, Name: doveSoap, Price: 39.99},
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *entity.Product
		wantErr error
	}{
		{
			name:    "product productService not initialized",
			fields:  fields{},
			args:    args{productID: doveSoap},
			wantErr: ProductNotFound,
		},
		{
			name: "product productService not initialized",
			fields: fields{
				master: productMaster,
			},
			args:    args{productID: "unknown product"},
			wantErr: ProductNotFound,
		},
		{
			name: "product found",
			fields: fields{
				master: productMaster,
			},
			args: args{productID: doveSoap},
			want: &entity.Product{
				ID:    doveSoap,
				Name:  doveSoap,
				Price: 39.99,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProductService(tt.fields.master)
			got, err := p.Lookup(tt.args.productID)
			if tt.wantErr != nil {
				require.NotNil(t, err, "error expected")
				assert.Equal(t, tt.wantErr.Error(), err.Error(), "error")
				return
			}

			require.NoError(t, err, "error")
			require.EqualValues(t, tt.want, got, "product found")
		})
	}
}
