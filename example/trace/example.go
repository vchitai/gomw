package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/vchitai/gomw"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var instrumentationName = "net/http"

type ctxKey string

func TraceStartHTTPMiddleware() gomw.HTTPMiddleware {
	return gomw.NewHTTPMiddleware(func(writer http.ResponseWriter, request *http.Request) (*http.Request, bool) {
		ctx, span := otel.GetTracerProvider().Tracer(instrumentationName).Start(request.Context(), "server request")
		ctx = context.WithValue(ctx, ctxKey("span-end"), span)
		request = request.WithContext(ctx)
		return request, true
	}, func(response gomw.HTTPResponse, request *http.Request) gomw.HTTPResponse {
		span, ok := request.Context().Value(ctxKey("span-end")).(trace.Span)
		if !ok {
			return response
		}
		span.End()
		return response
	})
}

func helloWorld() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello world!"))
	}
}

func main() {
	var mux = http.NewServeMux() // create new server

	mux.Handle("/", TraceStartHTTPMiddleware()(helloWorld()))

	l, err := net.Listen("tcp", ":10080")
	if err != nil {
		log.Fatal(err)
	}

	if err := http.Serve(l, mux); err != nil {
		log.Fatal(err)
	}
}
