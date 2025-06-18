package main

import (
	"fmt"
	"net"
)

func main() {
	var port string = ":8080"

	listener, err := net.Listen("tcp", port)

	if err != nil {
		fmt.Printf("Failed to listen on port %s\n", port)
		return
	}

	defer listener.Close()

	fmt.Printf("Server listening on port %s\n", port)

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Printf("Failed to accept connection: %v", err)
			continue
		}

		fmt.Println("Success!")
		conn.Close()
	}
}
