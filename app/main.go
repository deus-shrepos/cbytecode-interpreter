package main

import (
	"log"
	"net"
	"os"
	"strconv"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	logger := log.Default()
	router := HttpRouter{}

	echoEndpoint := func(r *HttpRequest, w *HttpResponse) {
		text := r.Params["str"]
		w.SetStatus(HttpStatusOk)
		w.SetHeader("Content-length", strconv.Itoa(len(text)))
		w.SetHeader("Content-Type", "text/plain")
		w.Write([]byte(text))
	}
	router.Handle("GET", "/echo/{str}", echoEndpoint)
	s, err := NewHttpServer("0.0.0.0", "4221", router, logger)
	if err != nil {
		panic(err)
	}

	s.ListenAndServe()
}
