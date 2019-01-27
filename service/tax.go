package service

type TaxService interface {
	ComputeTax(total float64) float64
}

type taxService struct {
	rate float64
}

func (t *taxService) ComputeTax(total float64) float64 {
	return total * t.rate / 100
}

func NewTaxService(rate float64) TaxService {
	return &taxService{rate}
}
