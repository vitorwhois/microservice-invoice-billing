package invoice

import "context"

type Repository interface {
	Create(ctx context.Context, invoice *Invoice) error
	GetByID(ctx context.Context, id int) (*Invoice, error)
	Update(ctx context.Context, invoice *Invoice) error
	List(ctx context.Context) ([]*Invoice, error)
	AddItem(ctx context.Context, item *InvoiceItem) error
}
