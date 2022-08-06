package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("Starting a gRPC server...")
	addr := fmt.Sprintf(":%d", 50051)
	Serve(addr)
}
