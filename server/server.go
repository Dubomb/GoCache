package main

import (
	"flag"
	"fmt"
	"gocache/cache"
	"net"
)

func main() {
	capacity := flag.Uint("capacity", 1000, "Maximum cache capacity")
	policy := flag.String("policy", "LRU", "Cache eviction policy: LRU or LFU")

	flag.Parse()

	var parsedPolicy cache.EvictionPolicy

	switch *policy {
	case "LRU":
		parsedPolicy = cache.NewLRUPolicy()
	case "LFU":
		parsedPolicy = cache.NewLFUPolicy()
	default:
		fmt.Println("Policy must be either LRU or LFU, please provide a valid policy")
		return
	}

	port := ":8080"

	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Printf("Failed to listen on port %s\n", port)
		return
	}

	defer listener.Close()

	fmt.Printf("Server listening on port %s\n", port)

	fmt.Printf("Initializing cache with capacity %d and policy %s!\n", *capacity, *policy)

	cache := cache.NewGoCache(*capacity, parsedPolicy)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleClientConnection(conn, cache)
	}
}
