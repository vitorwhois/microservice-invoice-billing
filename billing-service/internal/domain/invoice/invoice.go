package invoice

import (
	"errors"
	"log"
	"time"
)

var (
	ErrInvalidStatus = errors.New("invalid invoice status")
	ErrAlreadyClosed = errors.New("invoice already closed")
	ErrEmptyInvoice  = errors.New("invoice has no items")
	ErrNotFound      = errors.New("invoice not found")
)

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusClosed Status = "CLOSED"
)

type InvoiceItem struct {
	ID        int
	InvoiceID int
	ProductID int
	Quantity  int
	Price     float64
	Name      string
}

type Invoice struct {
	ID         int
	Number     string
	Status     Status
	CreatedAt  time.Time
	ClosedAt   *time.Time
	Items      []*InvoiceItem
	TotalValue float64
}

func NewInvoice(number string) *Invoice {
	return &Invoice{
		Number:    number,
		Status:    StatusOpen,
		CreatedAt: time.Now(),
		Items:     make([]*InvoiceItem, 0),
	}
}

func (i *Invoice) AddItem(productID int, quantity int, price float64, name string) *InvoiceItem {
	item := &InvoiceItem{
		ProductID: productID,
		Quantity:  quantity,
		Price:     price,
		Name:      name,
	}
	log.Printf("Adicionando item ao invoice: %v", item)

	i.Items = append(i.Items, item)

	i.CalculateTotal()
	log.Printf("Item %d adicionado ao invoice %d", item.ID, i.ID)
	return item
}

func (i *Invoice) Close() error {
	if i.Status == StatusClosed {
		return ErrAlreadyClosed
	}

	if len(i.Items) == 0 {
		return ErrEmptyInvoice
	}

	i.Status = StatusClosed
	now := time.Now()
	i.ClosedAt = &now

	return nil
}
func (i *Invoice) CalculateTotal() {
	total := 0.0
	for _, item := range i.Items {
		total += float64(item.Quantity) * item.Price
	}
	i.TotalValue = total
}
