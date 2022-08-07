package main

import (
	"fmt"

	"github.com/Monrevil/simplified-market-ledger/ledger"
)

func main() {
	addr := fmt.Sprintf(":%d", 50051)
	ledger.Serve(addr)
}
