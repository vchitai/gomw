package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/vchitai/gomw/common"
)

func main() {
	var mux = http.NewServeMux() // create new server
	mux.Handle("/", common.ValidationHTTPMiddleware(func(request *http.Request) (error) {
		bodyReader, err := request.GetBody()
		if err != nil {
			return err
		}
		var body []byte
		_, err = bodyReader.Read(body)
		if err != nil {
			return err
		}
		if string(body) != "hello" {
			return fmt.Errorf("not a valid request")
		}
		return nil
	})(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
