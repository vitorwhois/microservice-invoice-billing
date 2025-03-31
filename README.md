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
│
├── billing-service/
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       ├── domain/
│       │   ├── invoice/
│       │   │   ├── invoice.go
│       │   │   ├── item.go
│       │   │   ├── status.go
│       │   │   └── repository.go
│       │   └── errors/
│       ├── application/
│       │   ├── commands/
│       │   │   ├── create_invoice.go
│       │   │   └── print_invoice.go
│       │   ├── queries/
│       │   └── services/
│       ├── infrastructure/
│       │   ├── persistence/
│       │   └── http/
│       │       ├── handlers/
│       │       └── routes/
│       └── shared/
│
├──inventory-service/
│   ├──cmd/
│   └── main.go                
├──internal/
│   ├── config/
│   ├──domain/
│   │   └──product/             
│   │       ├── product.go       
│   │       └── repository.go    
│   ├── infrastructure/
│   │   ├── persistence/        
│   │   │   └── postgres.go
│   │   └── http/               
│   │       └── handlers.go
│   └── application/          
│       └── product/
│           └── service.go     
│
├── docker-compose.yml
└── README.md
```

## Tecnologias

- Golang
- PostgreSQL
- Docker e Docker Compose
- Gorilla Mux

## Requisitos

- Go 1.19+
- Docker e Docker Compose
- PostgreSQL 13+

### Melhorias Futuras
- Adicionar Circuit Breaker para chamadas entre serviços
- Mensageria para comunicação e filas com RabbitMq ou Kafka
- Implementar Dead Letter Queue para retry de falhas


### 1. Clonar repositório
```bash
git clone git@github.com:vitorwhois/billing-invoice-service-teste.git
cd billing-invoice-service-teste
```

### 2. Clonar repositório
```bash
docker-compose up --build 
```

