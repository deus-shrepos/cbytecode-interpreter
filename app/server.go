package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type handleFunc func(*HttpRequest, *HttpResponse)
type HttpServer struct {
	listener  net.Listener
	endpoints map[string]handleFunc
	log       *log.Logger
}

func NewHttpServer(addr string, port string, log *log.Logger) (*HttpServer, error) {
	listener, err := net.Listen("tcp",
		strings.Join([]string{addr, ":", port}, ""))
	if err != nil {
		return nil, err
	}
	return &HttpServer{
		listener:  listener,
		endpoints: make(map[string]handleFunc, 64),
		log:       log,
	}, nil
}

func (s *HttpServer) setEndpoint(path string, handler handleFunc) {
	s.endpoints[path] = handler
}

func (s *HttpServer) ListenAndServe() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.log.Fatalf("Could not accept new client request: %v", err)
		}

		go s.Handler(conn)
	}
}

func (s *HttpServer) Handler(conn net.Conn) {
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

	go func() {
		s.log.Println("Got request for path: ", request.Statusline.path)
		handler, exists := s.endpoints[request.Statusline.path]
		if !exists {
			response.SetStatus(HttpForbidden)
			response.Write([]byte(`{"error": "Not found"}`))
			return
		}
		s.log.Println("Running the handler now...")
		handler(&request, &response)
	}()
}
