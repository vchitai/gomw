package main

import (
	"log"
	"net"
	"net/http"

	"github.com/vchitai/gomw/common"
)

func main() {
	var mux = http.NewServeMux() // create new server
	mux.Handle("/", common.PayloadLoggingHTTPMiddleware()(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello world!"))
	})))

	l, err := net.Listen("tcp", ":10080")
	if err != nil {
		log.Fatal(err)
	}

	if err := http.Serve(l, mux); err != nil {
		log.Fatal(err)
	}
}
