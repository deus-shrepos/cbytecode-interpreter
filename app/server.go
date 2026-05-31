package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type HttpServer struct {
	listener net.Listener
	Router   HttpRouter
	log      *log.Logger
}

func NewHttpServer(addr string, port string, router HttpRouter, log *log.Logger) (*HttpServer, error) {
	listener, err := net.Listen("tcp",
		strings.Join([]string{addr, ":", port}, ""))
	if err != nil {
		return nil, err
	}
	return &HttpServer{
		listener: listener,
		Router:   router,
		log:      log,
	}, nil
}

func (s *HttpServer) ListenAndServe() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.log.Fatalf("Could not accept new client request: %v", err)
		}
		go s.HandleConn(conn)
	}
}

func (s *HttpServer) HandleConn(conn net.Conn) {
	buf := make([]byte, 1024)
	defer conn.Close()
	_, err := conn.Read(buf)
	if err != nil {
		s.log.Printf("Could not read into the buffer: %v", err)
	}
	request, err := ParseHttpRequest(buf)
	response := NewHttpResponse(conn)
	if err != nil {
		switch {
		case errors.Is(err, ErrMalformedRequest):
			response.SetStatus(HttpInvalidRequest)
		case errors.Is(err, ErrInvalidPath):
			response.SetStatus(HttpNotFound)
		default:
			response.SetStatus(HttpInternalServerError)
			response.Write([]byte(fmt.Sprintf(`{"error": %v}`, err)))
		}
		return
	}
	s.log.Printf("Path request by the client %v", request.StartLine.path)
	// We dispatch the request & response to the router
	s.Router.Dispatch(&request, &response)
}
