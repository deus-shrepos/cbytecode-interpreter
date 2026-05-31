package main

import (
	"bytes"
	"fmt"
	"net"
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

func (r *HttpResponse) Write(body []byte) (int, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "HTTP/1.1 %v %s\r\n\r\n", r.StatusLine.status, r.StatusLine.reason)
	for k, v := range r.Headers {
		fmt.Fprintf(&buf, "%v:%v\r\n", k, v)
	}
	fmt.Fprintf(&buf, "\r\n\r\n")
	buf.Write(body)
	return r.Conn.Write(buf.Bytes())
}

func (r *HttpResponse) SetHeader(key, value string) {
	r.Headers[key] = value
}

func (r *HttpResponse) SetStatus(statusCode StatusCode) {
	r.StatusLine.status = statusCode.Code
	r.StatusLine.reason = statusCode.Reason
	r.StatusLine.version = "HTTP/1.1" // static for now
}
