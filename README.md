# gomw - go-http-middleware

[![Godoc Reference](https://godoc.org/github.com/vchitai/zapflatencoder?status.svg)](http://godoc.org/github.com/vchitai/zapflatencoder)
[![Go Report Card](https://goreportcard.com/badge/github.com/vchitai/zap-flat-encoder)](https://goreportcard.com/report/github.com/vchitai/zap-flat-encoder)

A framework for fasten creating Go HTTP Middleware

## Get it

```shell
$ go get -u github.com/vchitai/gomw
```

## Usage

Implement a simple authenticate middleware with gomw, find it at [example](./example)

```go
package main

import (
	"log"
	"net"
	"net/http"

	"github.com/vchitai/gomw"
)

const AccessToken = "something"

func isValidToken(token string) bool {
	if token == AccessToken {
		return true
	}
	return false
}

func main() {
	authenticateMw := gomw.NewHTTPMiddleware(func(writer http.ResponseWriter, request *http.Request) (*http.Request, bool) {
		token := request.Header.Get("Authorization")
		if isValidToken(token) {
			return request, true
		}

		writer.WriteHeader(http.StatusUnauthorized)
		_, _ = writer.Write([]byte(http.StatusText(http.StatusUnauthorized)))
		return nil, false
	}, func(resp gomw.HTTPResponse) gomw.HTTPResponse {
		return gomw.NewHTTPResponse(resp.Body(), resp.Code()) // do some edit
	})

	var mux = http.NewServeMux() // create new server
	mux.Handle("/", authenticateMw(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
```
