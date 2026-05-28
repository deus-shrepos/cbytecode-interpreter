package main

import (
	"errors"
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		_ = fmt.Errorf("Error reading into the buffer: %v", err.Error())
	}
	// Read/Parse the Request header
	_, err = ParseHttpRequest(buf)
	response := HttpResponse{}
	if err != nil {
		switch {
		case errors.Is(err, ErrMalformedRequest):
			response.SetStatus(HttpInvalidRequest)
		case errors.Is(err, ErrInvalidPath):
			response.SetStatus(HttpNotFound)
		default:
			response.SetStatus(HttpInternalServerError)
		}
	} else {
		response.SetStatus(HttpStatusOk)
	}
	_, err = conn.Write(response.toBytes())
	if err != nil {
		_ = fmt.Errorf("Error writing back to the client: %v", err.Error())
	}
}
