# Simplified Market Ledger

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

# Process
1. Build basic in-memory solution. With an ability to sell/buy invoice. Define basic data structures.


# Decisions

## Single go module (go.mod file)
Pros:
- Simple way to update dependencies for all services
- No need to use to workspaces in vsCode
Cons:
- Longer build times for Docker images, if dependencies are no cached

## [Saga pattern](https://microservices.io/patterns/data/saga.html)
