package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/application/invoice"
	"github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/config"
	httphandlers "github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/infrastructure/http/handlers"
	"github.com/vitorwhois/microservice-invoice-billing/billing-service/internal/infrastructure/persistence"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	invoiceRepo := persistence.NewInvoiceRepository(db)
	invoiceService := invoice.NewInvoiceService(invoiceRepo, cfg.InventoryServiceURL)
	invoiceHandler := httphandlers.NewInvoiceHandler(invoiceService)

	router := mux.NewRouter()
	router.HandleFunc("/invoices", invoiceHandler.CreateInvoice).Methods("POST")
	router.HandleFunc("/invoices", invoiceHandler.ListInvoices).Methods("GET")
	router.HandleFunc("/invoices/{id}", invoiceHandler.GetInvoice).Methods("GET")
	router.HandleFunc("/invoices/{id}/items", invoiceHandler.AddInvoiceItem).Methods("POST")
	router.HandleFunc("/invoices/{id}/print", invoiceHandler.PrintInvoice).Methods("POST")

	router.Use(loggingMiddleware)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + getPort(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Billing service running on port %s", getPort())
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	return port
}

func setupDatabase(cfg *config.Config) (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Name)
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
