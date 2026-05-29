package main

import (
	"bytes"
	"net"
	"strconv"
)

type StatusLine struct {
	version string
	status  int
	reason  string
}
type HttpResponse struct {
	StatusLine StatusLine
	Headers    map[string]string
	Conn       net.Conn
}

func NewHttpResponse(conn net.Conn) HttpResponse {
	return HttpResponse{
		Conn:    conn,
		Headers: make(map[string]string, 256),
	}
}

func (r *HttpResponse) Write(b []byte) (int, error) {
	s := bytes.Buffer{}
	s.WriteString(r.StatusLine.version + " ")
	s.WriteString(strconv.Itoa(r.StatusLine.status) + " ")
	s.WriteString(r.StatusLine.reason)
	s.WriteString("\r\n\r\n")

	for k, v := range r.Headers {
		s.WriteString(k + ":" + v)
		s.WriteString("\r\n")
	}
	s.WriteString("\r\n")
	s.Write(b)
	return r.Conn.Write(s.Bytes())
}

func (r *HttpResponse) SetHeader(key, value string) {
	r.Headers[key] = value
}

func (r *HttpResponse) SetStatus(statusCode StatusCode) {
	r.StatusLine.status = statusCode.Code
	r.StatusLine.reason = statusCode.Reason
	r.StatusLine.version = "HTTP/1.1" // static for now
}
