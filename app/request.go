package main

import (
	"bytes"
	"net"
	"path"
	"regexp"
	"strings"
)

type RequestStartLine struct {
	method  string
	path    string
	version string
}

type HttpRequest struct {
	StartLine RequestStartLine
	Headers   map[string]string
	Body      []byte
	Params    map[string]string
}

var methodRe, _ = regexp.Compile(`POST|GET|OPTIONS|DELETE|PUT|UPDATE|PATCH`)

// we need to capture all four http path forms:
// origin-form: /index.html?x=1 HTTP/1.1
// absolute-form: https://example.com/index.html HTTP/1.1
// authority-form: example.com:443 HTTP/1.1
// asterik-form: * HTTP/1.1
var pathRe, _ = regexp.Compile(`(\*|https?:\/\/[^\s]+|\/[^\s]*|[A-Za-z0-9.-]+:\d+)`)
var versionRe, _ = regexp.Compile(`HTTP\/\d.\d`)

var placeholderRE = regexp.MustCompile(`\{([^{}]+)\}`)

func ParseHttpRequest(data []byte) (HttpRequest, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return HttpRequest{}, NewHttpError(
			ErrMalformedRequest,
			"CRLF missing.",
		)
	}
	startLine := bytes.SplitN(bytes.Trim(data[:idx], " "), []byte(" "), 3)
	if len(startLine) != 3 {
		return HttpRequest{}, NewHttpError(
			ErrMalformedRequest,
			"invalid start-line length.",
		)
	}
	method, urlPath, version := string(startLine[0]), string(startLine[1]), string(startLine[2])
	urlPath = path.Clean(urlPath)
	if !methodRe.MatchString(method) {
		return HttpRequest{}, NewHttpError(
			ErrMalformedRequest,
			"invalid method",
		)
	}

	if !pathRe.MatchString(urlPath) {
		return HttpRequest{}, NewHttpError(
			ErrMalformedRequest,
			"invalid path.",
		)
	}

	if !versionRe.MatchString(version) {
		return HttpRequest{}, NewHttpError(
			ErrMalformedRequest,
			"invalid version",
		)
	}

	if urlPath == "*" && method != "OPTIONS" {
		return HttpRequest{}, NewHttpError(
			ErrMalformedRequest,
			"'*' expects OPTIONS",
		)
	}

	if method == "CONNECT" {
		_, _, err := net.SplitHostPort(urlPath)
		if err != nil {
			return HttpRequest{}, NewHttpError(
				ErrMalformedRequest,
				"invalid host:port format",
			)
		}
	}

	requestStartLine := RequestStartLine{
		method:  method,
		path:    urlPath,
		version: version,
	}
	data = data[idx+2:]
	headers := make(map[string]string, 1024)
	for {
		line, rest, ok := bytes.Cut(data, []byte("\r\n\r\n"))
		if !ok {
			break
		}
		header := strings.SplitN(string(line), ":", 2)
		if len(header) != 2 {
			return HttpRequest{}, NewHttpError(
				ErrMalformedRequest,
				"invalid request header.",
			)
		}
		key, value := header[0], header[1]
		if !validateHeaderKey(key) {
			return HttpRequest{}, NewHttpError(
				ErrMalformedRequest,
				"invalid header key.",
			)
		}
		headers[key] = value
		data = rest
	}

	httpRequest := HttpRequest{
		StartLine: requestStartLine,
		Headers:   headers,
		Body:      data,
	}

	return httpRequest, nil
}

func validateHeaderKey(key string) bool {
	for _, c := range key {
		if ('a' <= c && c <= 'z') ||
			('A' <= c && c <= 'Z') ||
			('0' <= c && c <= '9') {
			continue
		}
		switch c {
		case '!', '$', '-', '_', '%', '*',
			'+', '^', '`', '|', '~', '#', '&', '\'':
			continue
		default:
			return false
		}
	}
	return true
}
