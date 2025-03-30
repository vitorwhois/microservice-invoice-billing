# Sistema de Faturamento e Notas Fiscais

Microsserviço para gestão de produtos e emissão de notas fiscais com controle de estoque.

## Escopo do Projeto

### Microserviço de Estoque

**Funcionalidades:**
- `POST /produtos`: Cadastro de produtos com:
  - ID
  - Nome
  - Preço
  - Saldo em estoque
- `POST /estoque/baixa`: Endpoint para baixa de estoque:
  - Valida saldo disponível
  - Atualiza quantidade em estoque
  - Simulação de falha via parâmetro `?fail=true`

### Microserviço de Faturamento

**Funcionalidades:**
- `POST /notas-fiscais`: Criação de notas fiscais com:
  - Número único
  - Status inicial (aberto)
  - Itens (produto e quantidade)
- `POST /notas-fiscais/{id}/imprimir`: Fluxo de impressão:
  1. Chama serviço de Estoque para validar saldos
  2. Executa baixa de estoque para todos os itens
  3. Atualiza status para "fechada" e registra data de emissão
  4. Retorna feedback detalhado ao usuário

**Comunicação entre Serviços:**
- Comunicação via HTTP REST
- Transacionalidade:
  - Rollback automático em caso de falha
  - Retentativas configuráveis
  - Feedback consistente de erros

## Estrutura Básica do Projeto
```bash
/project-root
  /estoque
    /cmd
      main.go
    /internal
      /domain      
      /handler    
      /repository  
  /faturamento
    /cmd
      main.go
    /internal
      /domain      
      /handler    
      /repository  
  docker-compose.yml 
```

## Tecnologias

- Golang
- PostgreSQL
- Docker e Docker Compose
- Gorilla Mux

## 📋 Requisitos

- Go 1.19+
- Docker e Docker Compose
- PostgreSQL 13+

## ⚙️ Configuração

### 1. Clonar repositório
```bash
git clone git@github.com:vitorwhois/billing-invoice-service-teste.git
cd billing-invoice-service-teste
```

### 2. Clonar repositório
```bash
docker-compose up --build 
```
