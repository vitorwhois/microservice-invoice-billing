package invoice

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
type RecoveryDetails struct {
	Attempted  bool     `json:"attempted"`
	Successful bool     `json:"successful"`
	Message    string   `json:"message,omitempty"`
	Details    []string `json:"details,omitempty"`
}

type InvoiceProcessResult struct {
	Success      bool            `json:"success"`
	InvoiceID    int             `json:"invoice_id"`
	StepReached  string          `json:"step_reached"`
	FailedReason string          `json:"failed_reason,omitempty"`
	Recovery     RecoveryDetails `json:"recovery,omitempty"`
}

var (
	ErrInventoryService  = errors.New("error communicating with inventory service")
	ErrProductNotFound   = errors.New("product not found")
	ErrStockReservation  = errors.New("failed to reserve stock")
	ErrStockConfirmation = errors.New("failed to confirm stock")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("invalid quantity")
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
		log.Printf("Fatura %d não existe", invoiceID)
		return err
	}

	if inv.Status == domaininvoice.StatusClosed {
		return domaininvoice.ErrAlreadyClosed
	}

	product, err := s.getProductFromInventory(ctx, productID, quantity)
	if err != nil {
		log.Printf("Produto %d nao encontrado", productID)
		return ErrProductNotFound
	}

	err = s.reserveStock(ctx, productID, quantity)
	if err != nil {
		log.Printf("Erro ao reservar estoque para o produto %d", productID)
		return ErrStockReservation
	}

	item := &domaininvoice.InvoiceItem{
		InvoiceID: invoiceID,
		ProductID: productID,
		Quantity:  quantity,
		Price:     product.Price,
		Name:      product.Name,
	}

	if err := s.repo.AddItem(ctx, item); err != nil {
		return err
	}

	inv, err = s.repo.GetByID(ctx, invoiceID)
	if err != nil {
		return err
	}

	log.Printf("Item %d adicionado ao invoice %d", productID, invoiceID)

	return s.repo.Update(ctx, inv)
}

func (s *Service) PrintInvoice(ctx context.Context, invoiceID int) (*InvoiceProcessResult, error) {
	inv, err := s.repo.GetByID(ctx, invoiceID)
	log.Printf("Verificando se a fatura %d existe", invoiceID)
	if err != nil {
		log.Printf("Fatura %d não existe", invoiceID)
		return nil, err
	}

	if inv.Status == domaininvoice.StatusClosed {
		return nil, domaininvoice.ErrAlreadyClosed
	}

	// Result object to track what happened during processing
	result := &InvoiceProcessResult{
		Success:      false,
		InvoiceID:    invoiceID,
		StepReached:  "started",
		FailedReason: "",
		Recovery:     RecoveryDetails{Attempted: false, Successful: false},
	}

	// Start transaction Saga
	reservedItems := make(map[int]int)

	// Step 1: Reserve stock for all items
	result.StepReached = "stock_reservation"
	for _, item := range inv.Items {
		if err := s.reserveStock(ctx, item.ProductID, item.Quantity); err != nil {
			// Compensating transaction: Cancel all reservations
			result.FailedReason = fmt.Sprintf("Failed to reserve stock for product %d: %v", item.ProductID, err)
			result.Recovery.Attempted = true
			result.Recovery.Message = "Attempting to cancel existing reservations"

			// Cancel previous reservations
			for prodID, qty := range reservedItems {
				result.Recovery.Details = append(result.Recovery.Details,
					fmt.Sprintf("Canceling reservation for product %d, quantity %d", prodID, qty))
				s.cancelReservation(ctx, prodID, qty)
			}

			result.Recovery.Successful = true
			return result, err
		}
		reservedItems[item.ProductID] = item.Quantity
	}

	// Step 2: Confirm all reservations
	result.StepReached = "stock_confirmation"
	for prodID, qty := range reservedItems {
		if err := s.confirmStock(ctx, prodID, qty); err != nil {
			// If confirming fails, cancel remaining reservations and try to restore confirmed ones
			result.FailedReason = fmt.Sprintf("Failed to confirm stock for product %d: %v", prodID, err)
			result.Recovery.Attempted = true
			result.Recovery.Message = "Attempting to cancel remaining reservations"

			log.Printf("Erro ao confirmar reserva de estoque para o produto %d", prodID, err)
			for pID, q := range reservedItems {
				if pID != prodID {
					result.Recovery.Details = append(result.Recovery.Details,
						fmt.Sprintf("Canceling reservation for product %d, quantity %d", pID, q))
					log.Printf("Cancelando reserva de estoque para o produto %d", pID)
					s.cancelReservation(ctx, pID, q)
				}
			}

			result.Recovery.Successful = true
			return result, err
		}
	}

	// Step 3: Cancel reservations before closing invoice
	result.StepReached = "cancel_stock_reservation"
	for prodID, qty := range reservedItems {
		log.Printf("Cancelando reserva de estoque para o produto %d, quantidade %d", prodID, qty)
		if err := s.cancelReservation(ctx, prodID, qty); err != nil {
			log.Printf("Failed to cancel reservation for product %d: %v", prodID, err)
			result.FailedReason = fmt.Sprintf("Error canceling reservation for product %d: %v", prodID, err)
			result.Recovery.Attempted = true
			result.Recovery.Details = append(result.Recovery.Details,
				fmt.Sprintf("Failed to cancel reservation for product %d", prodID))
			return result, err
		}
	}

	// 4: Close invoice
	result.StepReached = "invoice_closing"
	if err := inv.Close(); err != nil {
		result.FailedReason = fmt.Sprintf("Failed to close invoice: %v", err)
		return result, err
	}

	if err := s.repo.Update(ctx, inv); err != nil {
		result.FailedReason = fmt.Sprintf("Failed to update invoice in database: %v", err)
		return result, err
	}

	result.Success = true
	return result, nil
}

func (s *Service) getProductFromInventory(ctx context.Context, productID int, quantity int) (*ProductResponse, error) {
	url := fmt.Sprintf("%s/products/%d", s.inventoryServiceURL, productID)
	log.Printf("Buscando produto %d no inventário", productID)

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
		fmt.Println("Error decoding response from inventory service:", err)
		return nil, err
	}

	if product.Stock < quantity {
		return nil, ErrStockReservation
	}
	log.Printf("Produto encontrado: %v", product)
	return &product, nil
}

func (s *Service) reserveStock(ctx context.Context, productID int, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/reserve-stock", s.inventoryServiceURL, productID)
	requestBody, _ := json.Marshal(map[string]int{"quantity": quantity})

	log.Printf("Enviando requisição para reservar estoque: %s, Body: %s", url, string(requestBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Println("Erro ao criar requisição:", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Erro ao fazer requisição para reservar estoque:", err)
		return ErrInventoryService
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Falha ao reservar estoque. Status: %d, Response: %s", resp.StatusCode, string(body))
		return fmt.Errorf("falha ao reservar estoque, status: %d", resp.StatusCode)
	}

	log.Println("Estoque reservado com sucesso")
	return nil
}

func (s *Service) confirmStock(ctx context.Context, productID int, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/confirm-stock", s.inventoryServiceURL, productID)
	log.Printf("Confirmando reserva de estoque para o produto %d", productID)

	payload := map[string]int{"quantity": quantity}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error making request to inventory service:", err)
		return ErrInventoryService
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to confirm stock for product %d: %d", productID, resp.StatusCode)
		return ErrInventoryService
	}

	return nil
}

func (s *Service) cancelReservation(ctx context.Context, productID int, quantity int) error {
	url := fmt.Sprintf("%s/products/%d/cancel-reserve", s.inventoryServiceURL, productID)
	log.Printf("Cancelando reserva de estoque para o produto %d", productID)
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	payload := map[string]int{"quantity": quantity}
	jsonPayload, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error making request cancelReservation to inventory service:", err)
		return ErrInventoryService
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to cancel reservation for product %d: %d", productID, resp.StatusCode)
	}

	return nil
}
