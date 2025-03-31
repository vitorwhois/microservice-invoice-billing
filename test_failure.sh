#!/bin/bash

# Criar um produto
echo "Criando produto..."
curl -X POST http://localhost:8080/products \
    -H "Content-Type: application/json" \
    -d '{"name":"Test Product", "price":100, "stock":10}'
echo -e "\n"

# Criar uma fatura
echo "Criando invoice..."
curl -X POST http://localhost:8081/invoices \
    -H "Content-Type: application/json" \
    -d '{"number":"INV-001"}'
echo -e "\n"

# Adicionar um item na invoice
echo "Adicionando item na invoice..."
curl -X POST http://localhost:8081/invoices/1/items \
    -H "Content-Type: application/json" \
    -d '{"product_id":1, "quantity":5}'
echo -e "\n"

# Tentar imprimir a invoice 
echo "Tentando imprimir invoice (esperado falhar)..."
curl -X POST http://localhost:8081/invoices/1/print
echo -e "\n"

# Verificar o status da invoice após a falha
echo "Verificando status da invoice..."
curl -X GET http://localhost:8081/invoices/1
echo -e "\n"

# Checar o estoque do produto para garantir que ele foi restaurado
echo "Verificando estoque do produto..."
curl -X GET http://localhost:8080/products/1
echo -e "\n"
