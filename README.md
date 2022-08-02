# Simplified Market Ledger

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
