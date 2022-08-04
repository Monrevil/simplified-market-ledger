# Simplified Market Ledger

# How to run

```
git clone https://github.com/Monrevil/simplified-market-ledger
make docker
docker compose up
```

Download [bloomrpc](https://github.com/bloomrpc/bloomrpc)
```
brew install --cask bloomrpc
```
1. Press Import protos (Green + sign) 
2. Chose proto file at `/api/api.proto/`
3. SellInvoice
4. ListInvestors
5. PlaceBid
6. ApproveFinancing
7. ListInvestors to check if Investor obtained the invoice, and his balance has changed

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
Before that is should check if Invoice is Available for financing, and if bid is valid.

# Considirations:
So far all of the processes should be a single transaction in a database.
PlaceBid\Approve\Reverse should Commit or Rollback all write operations at the end.

# Q? 
1. Ledger should reflect bids that where smaller than an invoice value in database? They should be recorded, with status rejected.

# Decisions

## Single go module (go.mod file)
Pros:
- Simple way to update dependencies for all services
- No need to use to workspaces in vsCode
Cons:
- Longer build times for Docker images, if dependencies are no cached

## [Saga pattern](https://microservices.io/patterns/data/saga.html)
