# Simplified Market Ledger
Tech challenge for AREX Markets hiring process.

# Table of Contents
- [Simplified Market Ledger](#simplified-market-ledger)
- [Table of Contents](#table-of-contents)
- [How to run](#how-to-run)
  - [Copy this repository](#copy-this-repository)
  - [Run docker compose up](#run-docker-compose-up)
- [Test Workflow](#test-workflow)
  - [With bloomrpc](#with-bloomrpc)
  - [Steps](#steps)
- [Build project locally](#build-project-locally)
  - [To build docker image run](#to-build-docker-image-run)
  - [To regenerate gRPC](#to-regenerate-grpc)
- [Required endpoints](#required-endpoints)
- [Matching algorithm](#matching-algorithm)

# How to run

## Copy this repository
```
git clone https://github.com/Monrevil/simplified-market-ledger
```

## Run docker compose up
```
docker compose up
```
It will pull and run:
- image with this project from ghcr.io
- postgres image
  
# Test Workflow

## With bloomrpc
Download [bloomrpc](https://github.com/bloomrpc/bloomrpc)
```
brew install --cask bloomrpc
```
## Steps
1. Press Import protos (Green + sign) 
2. Chose proto file at `/api/api.proto/`. Ledger service will be running on port 5050
3. NewIssuer - create new issuer, get issuerID
4. SellInvoice - sell invoice, using issuerID, get invoiceID
5. NewInvestor - create new investor, get investorID
6. ListInvestors - check created investors
7. PlaceBid - place bid using investorID, and invoiceID, get transactionID
8. ApproveFinancing - approve financing, using transactionID
9. ListInvestors to check if Investor has obtained the invoice, and his balance has changed

# Build project locally
If you wish to build and run project on your machine, instead of pulling docker image, run:

```
go mod download -x
go run .
```

LedgerService will try to connect to postgres database at localhost:5432
If `POSTGRES_HOST` or `POSTGRES_PORT` env variables are set, it will use them instead of defaults.
## To build docker image run
```
make docker
```
## To regenerate gRPC
```
make proto
```

# Required endpoints
Endpoints defined in the documentation explicitly:
- An endpoint to sell an invoice
- An endpoint to retrieve an invoice, including the status: if it has been bought or not. In
case it has been bought, by which investor
- An endpoint to list all the investors and their balances

Derived from defined functionality:
- An endpoint to place a bid
- An endpoint to approve financing
- An endpoint to revert financing

# Matching algorithm
Matching algorithm should be Singleton, and multiplex all connections (bid attempts) for a given invoice into a single goroutine. 
It should check if Invoice is Available for financing, and if bid is valid.
A go channel is used to enforce FIFO order.