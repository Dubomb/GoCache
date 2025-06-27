package main

import (
	"bufio"
	"fmt"
	"gocache/parser"
	"io"
	"net"
)

func handleClientConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connection accepted from: ", conn.RemoteAddr())

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection to ", conn.RemoteAddr().String(), " terminated!")
			} else {
				fmt.Println("Error when reading: ", err)
			}

			return
		}

		command, err := parser.ParseCommand(msg)

		if err == nil {
			fmt.Println("Type: ", command.Type, " Args: ", command.Args)
		} else {
			fmt.Println("Error encountered: ", err)
		}

	}
}
