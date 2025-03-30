package product

import (
	"errors"
	"time"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidStock      = errors.New("invalid stock quantity")
	ErrNotFound          = errors.New("product not found")
	ErrConcurrentUpdate  = errors.New("concurrent modification")
)

type Product struct {
	ID            int
	Name          string
	Price         float64
	Stock         int
	ReservedStock int
	Version       int
	CreatedAt     time.Time
}

func NewProduct(name string, price float64, stock int) (*Product, error) {
	if stock < 0 {
		return nil, ErrInvalidStock
	}

	return &Product{
		Name:      name,
		Price:     price,
		Stock:     stock,
		Version:   1,
		CreatedAt: time.Now(),
	}, nil
}

func (p *Product) ReserveStock(quantity int) error {
	availableStock := p.Stock - p.ReservedStock
	if quantity > availableStock {
		return ErrInsufficientStock
	}
	p.ReservedStock += quantity
	p.Version++
	return nil
}

func (p *Product) ConfirmReservation(quantity int) error {
	if quantity > p.ReservedStock {
		return ErrInvalidStock
	}
	p.Stock -= quantity
	p.ReservedStock -= quantity
	p.Version++
	return nil
}

func (p *Product) CancelReservation(quantity int) error {
	if quantity > p.ReservedStock {
		return ErrInvalidStock
	}
	p.ReservedStock -= quantity
	p.Version++
	return nil
}
