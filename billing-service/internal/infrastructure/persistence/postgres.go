package persistence

import (
	"context"
	"database/sql"

	"github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/domain/invoice"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewInvoiceRepository(db *sql.DB) invoice.Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, inv *invoice.Invoice) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
        INSERT INTO invoices (number, status, created_at, total_value)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	err = tx.QueryRowContext(ctx, query,
		inv.Number, inv.Status, inv.CreatedAt, inv.TotalValue,
	).Scan(&inv.ID)

	if err != nil {
		return err
	}

	for _, item := range inv.Items {
		item.InvoiceID = inv.ID
		query := `
            INSERT INTO invoice_items (invoice_id, product_id, quantity, price, name)
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id`

		err = tx.QueryRowContext(ctx, query,
			item.InvoiceID, item.ProductID, item.Quantity, item.Price, item.Name,
		).Scan(&item.ID)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostgresRepository) GetByID(ctx context.Context, id int) (*invoice.Invoice, error) {
	invQuery := `
        SELECT id, number, status, created_at, closed_at, total_value
        FROM invoices 
        WHERE id = $1`

	inv := &invoice.Invoice{}
	err := r.db.QueryRowContext(ctx, invQuery, id).Scan(
		&inv.ID, &inv.Number, &inv.Status, &inv.CreatedAt, &inv.ClosedAt, &inv.TotalValue,
	)

	if err == sql.ErrNoRows {
		return nil, invoice.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	itemsQuery := `
        SELECT id, product_id, quantity, price, name
        FROM invoice_items
        WHERE invoice_id = $1`

	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	inv.Items = make([]*invoice.InvoiceItem, 0)
	for rows.Next() {
		item := &invoice.InvoiceItem{InvoiceID: id}
		if err := rows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price, &item.Name); err != nil {
			return nil, err
		}
		inv.Items = append(inv.Items, item)
	}
	inv.CalculateTotal()
	return inv, nil
}

func (r *PostgresRepository) Update(ctx context.Context, inv *invoice.Invoice) error {
	query := `
        UPDATE invoices
        SET status = $1, closed_at = $2, total_value = $3
        WHERE id = $4`

	_, err := r.db.ExecContext(ctx, query,
		inv.Status, inv.ClosedAt, inv.TotalValue, inv.ID,
	)

	return err
}
func (r *PostgresRepository) AddItem(ctx context.Context, item *invoice.InvoiceItem) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	itemQuery := `
        INSERT INTO invoice_items (invoice_id, product_id, quantity, price, name)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`

	err = tx.QueryRowContext(ctx, itemQuery,
		item.InvoiceID, item.ProductID, item.Quantity, item.Price, item.Name,
	).Scan(&item.ID)

	if err != nil {
		return err
	}

	updateQuery := `
        UPDATE invoices 
		SET total_value = total_value + CAST($1 AS numeric) * CAST($2 AS numeric)
        WHERE id = $3`

	_, err = tx.ExecContext(ctx, updateQuery,
		item.Quantity, item.Price, item.InvoiceID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PostgresRepository) List(ctx context.Context) ([]*invoice.Invoice, error) {
	query := `
        SELECT id, number, status, created_at, closed_at, total_value
        FROM invoices
        ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invoices := make([]*invoice.Invoice, 0)
	for rows.Next() {
		inv := &invoice.Invoice{}
		if err := rows.Scan(&inv.ID, &inv.Number, &inv.Status, &inv.CreatedAt, &inv.ClosedAt, &inv.TotalValue); err != nil {
			return nil, err
		}

		itemsQuery := `
			SELECT id, product_id, quantity, price, name
			FROM invoice_items
			WHERE invoice_id = $1`

		itemRows, err := r.db.QueryContext(ctx, itemsQuery, inv.ID)
		if err != nil {
			return nil, err
		}

		inv.Items = make([]*invoice.InvoiceItem, 0)
		for itemRows.Next() {
			item := &invoice.InvoiceItem{InvoiceID: inv.ID}
			if err := itemRows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.Price, &item.Name); err != nil {
				itemRows.Close()
				return nil, err
			}
			inv.Items = append(inv.Items, item)
		}
		itemRows.Close()

		invoices = append(invoices, inv)
	}

	return invoices, nil
}
