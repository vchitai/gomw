package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/vchitai/gomw/common"
)

type claim struct {
	UserID       string `json:"user_id"`
	UserPassword string `json:"user_password"`
}

type authClaimKey struct{}

func claimFromToken(token string) (*claim, error) {
	var decodedToken []byte
	if _, err := base64.NewDecoder(base64.StdEncoding, bytes.NewBufferString(token)).Read(decodedToken); err != nil {
		return nil, err
	}

	segments := strings.Split(string(decodedToken), ":")
	if len(segments) != 2 {
		return nil, fmt.Errorf("token does not have valid number of segments")
	}
	return &claim{
		UserID:       segments[0],
		UserPassword: segments[1],
	}, nil
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
