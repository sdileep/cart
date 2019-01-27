package entity

type Product struct {
	ID    string
	Name  string
	Price float64 // float64 for simplicity
}

type Cart struct {
	ID    string
	Items []*CartItem
	Total float64
}

type CartItem struct {
	ProductID string
	Quantity  uint8
	UnitPrice float64
}
