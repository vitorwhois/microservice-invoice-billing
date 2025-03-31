.PHONY: build run test clean docker-up docker-down


INVENTORY_SERVICE=./inventory-service
BILLING_SERVICE=./billing-service

build:
	@echo "Building services..."
	cd $(INVENTORY_SERVICE) && go build -o bin/inventory-service ./cmd/main.go
	cd $(BILLING_SERVICE) && go build -o bin/billing-service ./cmd/main.go

run-inventory:
	cd $(INVENTORY_SERVICE) && go run ./cmd/main.go

run-billing:
	cd $(BILLING_SERVICE) && go run ./cmd/main.go

test:
	cd $(INVENTORY_SERVICE) && go test ./...
	cd $(BILLING_SERVICE) && go test ./...

clean:
	rm -rf $(INVENTORY_SERVICE)/bin
	rm -rf $(BILLING_SERVICE)/bin

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose build

docker-logs:
	docker-compose logs -f

init-test-data:
	@echo "Inserindo dados de teste..."
	curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Notebook", "price":2800.00, "stock":10}'
	curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Mouse", "price":50.00, "stock":30}'
	curl -X POST http://localhost:8080/products -H "Content-Type: application/json" -d '{"name":"Teclado", "price":100.00, "stock":20}'
	curl -X POST http://localhost:8081/invoices -H "Content-Type: application/json" -d '{"number":"NF001"}'