package service

import (
	"math"
	"testing"
)

func Test_taxService_ComputeTax(t *testing.T) {
	type fields struct {
		rate float64
	}
	type args struct {
		total float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "sample",
			fields: fields{
				rate: 12.5,
			},
			args: args{
				total: 2*39.99 + 2*99.99,
			},
			want: 35.00,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewTaxService(tt.fields.rate)
			if got := service.ComputeTax(tt.args.total); math.Ceil(got*100)/100 != tt.want {
				t.Errorf("taxService.ComputeTax() = %v, want %v", got, tt.want)
			}
		})
	}
}
