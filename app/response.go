package main

import (
	"bytes"
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
}

func (r *HttpResponse) toBytes() []byte {
	s := bytes.Buffer{}
	s.WriteString(r.StatusLine.version + " ")
	s.WriteString(strconv.Itoa(r.StatusLine.status) + " ")
	s.WriteString(r.StatusLine.reason)
	s.WriteString("\r\n\r\n")

	for k, v := range r.Headers {
		s.WriteString(k + ":" + v)
		s.WriteString("\r\n")
	}
	return s.Bytes()
}

func (r *HttpResponse) SetHeader(key, value string) {
	r.Headers[key] = value
}

func (r *HttpResponse) SetStatus(statusCode StatusCode) {
	r.StatusLine.status = statusCode.Code
	r.StatusLine.reason = statusCode.Reason
	r.StatusLine.version = "HTTP/1.1" // static for now
}
