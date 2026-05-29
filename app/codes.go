package main

type StatusCode struct {
	Code   int
	Reason string
}

var HttpNotFound = StatusCode{Code: 404, Reason: "Not Found"}
var HttpInvalidRequest = StatusCode{Code: 400, Reason: "Bad Request"}
var HttpStatusOk = StatusCode{Code: 200, Reason: "OK"}
var HttpInternalServerError = StatusCode{Code: 500, Reason: "Internal Server Error"}
var HttpForbidden = StatusCode{Code: 403, Reason: "Forbidden"}
