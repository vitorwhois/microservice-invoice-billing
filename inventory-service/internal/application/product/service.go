package product

import (
	"context"

	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/domain/product"
)

type Service struct {
	repo product.Repository
}

func NewProductService(repo product.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateProduct(ctx context.Context, name string, price float64, stock int) (*product.Product, error) {
	product, err := product.NewProduct(name, price, stock)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, product)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *Service) ReserveStock(ctx context.Context, id int, quantity int) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := product.ReserveStock(quantity); err != nil {
		return err
	}

	return s.repo.Update(ctx, product)
}

func (s *Service) ConfirmStock(ctx context.Context, id int, quantity int) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := product.ConfirmReservation(quantity); err != nil {
		return err
	}

	return s.repo.Update(ctx, product)
}

func (s *Service) CancelReservation(ctx context.Context, id int, quantity int) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err := product.CancelReservation(quantity); err != nil {
		return err
	}

	return s.repo.Update(ctx, product)
}

func (s *Service) GetProductByID(ctx context.Context, id int) (*product.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}
