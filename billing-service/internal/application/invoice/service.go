package invoice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	domaininvoice "github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/domain/invoice"
)

type Service struct {
	repo                domaininvoice.Repository
	inventoryServiceURL string
}

type ProductResponse struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

var (
	ErrInventoryService  = errors.New("error communicating with inventory service")
	ErrProductNotFound   = errors.New("product not found")
	ErrStockReservation  = errors.New("failed to reserve stock")
	ErrStockConfirmation = errors.New("failed to confirm stock")
)

func NewInvoiceService(repo domaininvoice.Repository, inventoryURL string) *Service {
	return &Service{
		repo:                repo,
		inventoryServiceURL: inventoryURL,
	}
}

func (s *Service) CreateInvoice(ctx context.Context, number string) (*domaininvoice.Invoice, error) {
	inv := domaininvoice.NewInvoice(number)
	if err := s.repo.Create(ctx, inv); err != nil {
		return nil, err
	}
	return inv, nil
}

func (s *Service) GetInvoiceByID(ctx context.Context, id int) (*domaininvoice.Invoice, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListInvoices(ctx context.Context) ([]*domaininvoice.Invoice, error) {
	return s.repo.List(ctx)
}

func (s *Service) AddInvoiceItem(ctx context.Context, invoiceID int, productID int, quantity int) error {
	inv, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if inv.Status == domaininvoice.StatusClosed {
		return domaininvoice.ErrAlreadyClosed
	}

	product, err := s.getProductFromInventory(ctx, productID)
	if err != nil {
		return err
	}

	inv.AddItem(productID, quantity, product.Price, product.Name)

	return s.repo.Update(ctx, inv)
}

func (s *Service) PrintInvoice(ctx context.Context, invoiceID int) error {
	inv, err := s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	if inv.Status == domaininvoice.StatusClosed {
		return domaininvoice.ErrAlreadyClosed
	}

	// Start transaction Saga
	reservedItems := make(map[int]int)

	// Step 1: Reserve stock for all items
	for _, item := range inv.Items {
		if err := s.reserveStock(ctx, item.ProductID, item.Quantity); err != nil {
			// Compensating transaction: Cancel all reservations
			for prodID, qty := range reservedItems {
				s.cancelReservation(ctx, prodID, qty)
			}
			return err
		}
		reservedItems[item.ProductID] = item.Quantity
	}

	// Step 2: Confirm all reservations
	for prodID, qty := range reservedItems {
		if err := s.confirmStock(ctx, prodID, qty); err != nil {
			// If confirming fails, cancel remaining reservations and try to restore confirmed ones
			for pID, q := range reservedItems {
				if pID != prodID {
					s.cancelReservation(ctx, pID, q)
				}
			}
			return err
		}
	}

	// 3: Close invoice
	if err := inv.Close(); err != nil {
		return err
	}

	return s.repo.Update(ctx, inv)
}

func (s *Service) getProductFromInventory(ctx context.Context, productID int) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/products/%d", s.inventoryServiceURL, productID)
	log.Printf("Tentando acessar: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request to inventory service:", err)
		return nil, ErrInventoryService
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrProductNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrInventoryService
	}

	var product ProductResponse
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *Service) reserveStock(ctx context.Context, productID int, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/reserve", s.inventoryServiceURL, productID)

	payload := map[string]int{"quantity": quantity}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return ErrInventoryService
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResponse map[string]string
		json.NewDecoder(resp.Body).Decode(&errResponse)

		if resp.StatusCode == http.StatusConflict {
			return ErrStockReservation
		}
		return ErrInventoryService
	}

	return nil
}

func (s *Service) confirmStock(ctx context.Context, productID int, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/confirm", s.inventoryServiceURL, productID)

	payload := map[string]int{"quantity": quantity}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return ErrInventoryService
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrInventoryService
	}

	return nil
}

func (s *Service) cancelReservation(ctx context.Context, productID int, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/cancel", s.inventoryServiceURL, productID)

	payload := map[string]int{"quantity": quantity}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return ErrInventoryService
	}
	defer resp.Body.Close()

	// Even if cancellation fails, we just log it
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to cancel reservation for product %d: %d", productID, resp.StatusCode)
	}

	return nil
}
