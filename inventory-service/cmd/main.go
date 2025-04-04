package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/application/product"
	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/config"
	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/infrastructure/http/handlers"
	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/infrastructure/http/routes"
	"github.com/vitorwhois/microservice-invoice-billing/inventory-service/internal/infrastructure/persistence"

	_ "github.com/lib/pq"

	"github.com/rs/cors"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	db, err := setupDatabase(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados: %v", err)
	}
	defer db.Close()

	failureMode := os.Getenv("INVENTORY_FAILURE_MODE")
	if failureMode != "" {
		log.Printf("Running with failure mode: %s", failureMode)
	}

	productRepo := persistence.NewProductRepository(db)
	productService := product.NewProductService(productRepo, failureMode)
	productHandler := handlers.NewProductHandler(productService)

	router := routes.NewRouter(productHandler)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	handler := c.Handler(router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	log.Printf("Servidor rodando na porta %s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Falha no servidor: %v", err)
	}
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
