package product

import "context"

type Repository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id int) (*Product, error)
	Update(ctx context.Context, product *Product) error
	GetAll(ctx context.Context) ([]*Product, error)
}
