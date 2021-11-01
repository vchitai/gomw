package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/vchitai/gomw/common"
)

type claim struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
}

type authClaimKey struct{}

func claimFromToken(token string) (*claim, error) {
	var claim claim
	if err := json.Unmarshal([]byte(token), &claim); err != nil {
		return nil, fmt.Errorf("decode token error %w", err)
	}
	return &claim, nil
}

func ExtractClaimFromToken(request *http.Request) (*http.Request, error) {
	token := request.Header.Get("Authorization")
	if len(token) == 0 {
		return request, fmt.Errorf("require authorization token included")
	}

	claim, err := claimFromToken(token)
	if err != nil {
		return request, fmt.Errorf("extract claim error %w", err)
	}
	newCtx := context.WithValue(request.Context(), authClaimKey{}, claim)
	request = request.WithContext(newCtx)
	return request, nil
}

func main() {
	var mux = http.NewServeMux() // create new server
	mux.Handle("/", common.AuthenticationHTTPMiddleware(ExtractClaimFromToken)(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
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
