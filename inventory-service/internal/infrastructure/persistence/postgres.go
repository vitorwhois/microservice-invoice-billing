package persistence

import (
	"context"
	"database/sql"

	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/domain/product"
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

func (r *PostgresRepository) GetAll(ctx context.Context) ([]*product.Product, error) {
	query := `
        SELECT id, name, price, stock, reserved_stock, version, created_at
        FROM products
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]*product.Product, 0)
	for rows.Next() {
		p := &product.Product{}
		if err := rows.Scan(
			&p.ID, &p.Name, &p.Price, &p.Stock, &p.ReservedStock, &p.Version, &p.CreatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
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
