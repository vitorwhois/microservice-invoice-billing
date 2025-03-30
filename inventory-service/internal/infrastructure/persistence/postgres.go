package persistence

import (
	"context"
	"database/sql"
	"inventory-service/internal/domain/product"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) product.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, p *product.Product) error {
	query := `
        INSERT INTO products (name, price, stock, reserved_stock, version)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		p.Name, p.Price, p.Stock, p.ReservedStock, p.Version,
	).Scan(&p.ID)
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*product.Product, error) {
	p := &product.Product{}
	query := `
        SELECT id, name, price, stock, reserved_stock, version, created_at
        FROM products WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Price, &p.Stock, &p.ReservedStock, &p.Version, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, product.ErrNotFound
	}
	return p, err
}

func (r *PostgresRepository) Update(ctx context.Context, p *product.Product) error {
	query := `
        UPDATE products 
        SET stock = $1, reserved_stock = $2, version = $3
        WHERE id = $4 AND version = $5`

	result, err := r.db.ExecContext(ctx, query,
		p.Stock, p.ReservedStock, p.Version,
		p.ID, p.Version-1,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return product.ErrConcurrentUpdate
	}

	return nil
}
