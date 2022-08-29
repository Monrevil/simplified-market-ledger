# Simplified Market Ledger

Tech challenge for AREX Markets hiring process.

I feel like there is much to improve in this implementation. Although I am not sure in what direction I should go.
And at the moment it represents my technical skills more or less accurately.

Some considerations I have at the moment:

## Drop gorm, and use raw SQL instead

Orms can be convenient, but sometimes they can act in unpredictable way.
And hard to use if you need union/inner joins, complex queries.

## Unit test

Postgres repository probably should be invoked through interface. So that ledger can be tested without the database layer.
Initially it was, but I had some complications while trying to have transactional queries and interface in front of it.
And having to deal with in-memory implementation of that seemed like a lot of meaningless extra work.

Maybe I am wrong about that. Would appreciate feedback if this is the case.

## Matching algorithm FIFO order

I did find this part of requirements a bit confusing. As usually letting postgres mvcc handle order is good enough.
Or if you actually need FIFO some kind of message broker is used, like RabbitMQ.
Maybe I should have asked more questions about it, but also I felt like this was just as tech challenge.

Simplest way to enforce it for PlaceBid endpoint should look like this (pseudocode):

```go
type Server struct {
 requestQueue chan req
}

type req struct {
 InvoiceID int
}

// Launch infinite loop to process all calls sent to the queue
func StartMatchingAlgorithm(ch chan req) {
 for {
  request := <-ch
  // do work
 }
}

// Redirect all calls to the queue
func (s Server) PlaceBid(r req) {
 s.requestQueue <- r
}

func main() {
 // make a buffered channel to serve as queue
 // 500 is just a first number that came to mind, it should be enough to process calls most of the time
 // otherwise should be stress tested and benchmarked
 ch := make(chan req, 500)
 go StartMatchingAlgorithm(ch)
 s := Server{
  requestQueue: ch,
 }
 s.Serve()
}
```

This solution will put all calls from grpc into a single goroutine.
Alternatively - to have individual queues for each unique invoice - a map with all queues is required, that will be locked by a mutex. And there will be a lot of allocations, garbage collection, mutex locking and unlocking. And that drives complexity of a solution quite a bit, while not necessarily leading to a better performance. Or, at least this is my understating at the moment.

## How to run

### Copy this repository

```bash
git clone https://github.com/Monrevil/simplified-market-ledger

cd simplified-market-ledger
```

### Run docker compose up

```bash
docker compose up
```

Wait for the app to start up.
It will pull and run:

- image with this project from ghcr.io
- postgres image
  
## Test Workflow

Some general directions on how to test the code.

### With bloomrpc

Download [bloomrpc](https://github.com/bloomrpc/bloomrpc)

```bash
brew install --cask bloomrpc
```

### Steps

1. Press Import protos (Green + sign)
2. Chose proto file at `/api/api.proto` Ledger service will be running on port 5050
3. NewIssuer - create new issuer, get issuerID
4. SellInvoice - sell invoice, using issuerID, get invoiceID
5. NewInvestor - create new investor, get investorID
6. ListInvestors - check created investors
7. PlaceBid - place bid using investorID, and invoiceID, get transactionID
8. ApproveFinancing - approve financing, using transactionID
9. ListInvestors to check if Investor has obtained the invoice, and his balance has changed

## Build project locally

If you wish to build and run project on your machine, instead of pulling docker image, run:

```bash
go mod download -x
go run .
```

LedgerService will try to connect to postgres database at localhost:5432
If `POSTGRES_HOST` or `POSTGRES_PORT` env variables are set, it will use them instead of defaults.

### To build docker image run

```bash
make docker
```

### To regenerate gRPC

```bash
make proto
```

## Required endpoints

Endpoints defined in the documentation explicitly:

- An endpoint to sell an invoice
- An endpoint to retrieve an invoice, including the status: if it has been bought or not. In
case it has been bought, by which investor
- An endpoint to list all the investors and their balances

Derived from defined functionality:

- An endpoint to place a bid
- An endpoint to approve financing
- An endpoint to revert financing
