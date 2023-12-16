package main

import (
	"sync"
)

func main() {
	var wg sync.WaitGroup

	// Start the REST API
	wg.Add(1)
	go func() {
		defer wg.Done()
		StartRESTServer()
	}()

	// Start the gRPC Server
	go func() {
		defer wg.Done()
		InitGRPCServer()
	}()

	wg.Wait()
}
