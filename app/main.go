package main

import (
	"log"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports above (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	logger := log.Default()
	s, err := NewHttpServer("0.0.0.0", "4221", logger)
	if err != nil {
		panic(err)
	}

	echoEndpoint := func(r *HttpRequest, w *HttpResponse) {
		placeHolders := r.ExtractPlaceHolder()
		w.SetHeader("Content-type", "text/plain")
		w.SetStatus(HttpStatusOk)
		w.Write([]byte(placeHolders[0]))
	}

	s.setEndpoint("/echo/{str}", echoEndpoint)
	s.ListenAndServe()
}
