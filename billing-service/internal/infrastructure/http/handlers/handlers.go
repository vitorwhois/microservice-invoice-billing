package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	appinvoice "github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/application/invoice"
	domaininvoice "github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/domain/invoice"

	"github.com/gorilla/mux"
)

type InvoiceHandler struct {
	service *appinvoice.Service
}

func NewInvoiceHandler(service *appinvoice.Service) *InvoiceHandler {
	return &InvoiceHandler{service: service}
}

func (h *InvoiceHandler) CreateInvoice(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Number string `json:"number"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inv, err := h.service.CreateInvoice(r.Context(), request.Number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(inv)
}

func (h *InvoiceHandler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	inv, err := h.service.GetInvoiceByID(r.Context(), id)
	if err != nil {
		if err == domaininvoice.ErrNotFound {
			http.Error(w, "Invoice not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func (h *InvoiceHandler) ListInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.service.ListInvoices(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

func (h *InvoiceHandler) AddInvoiceItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	var request struct {
		ProductID int `json:"product_id"`
		Quantity  int `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.AddInvoiceItem(r.Context(), id, request.ProductID, request.Quantity)
	if err != nil {
		if err == domaininvoice.ErrAlreadyClosed {
			http.Error(w, "Invoice is already closed", http.StatusConflict)
		} else if err == appinvoice.ErrProductNotFound {
			http.Error(w, "Product not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return the updated invoice
	inv, _ := h.service.GetInvoiceByID(r.Context(), id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(inv)
}

func (h *InvoiceHandler) PrintInvoice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid invoice ID", http.StatusBadRequest)
		return
	}

	err = h.service.PrintInvoice(r.Context(), id)
	if err != nil {
		if err == domaininvoice.ErrAlreadyClosed {
			http.Error(w, "Invoice is already closed", http.StatusConflict)
		} else if err == appinvoice.ErrStockReservation {
			http.Error(w, "Insufficient stock for one or more products", http.StatusConflict)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Return success with the updated invoice
	inv, _ := h.service.GetInvoiceByID(r.Context(), id)

	response := struct {
		Message string                 `json:"message"`
		Invoice *domaininvoice.Invoice `json:"invoice"`
	}{
		Message: "Invoice printed successfully",
		Invoice: inv,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
