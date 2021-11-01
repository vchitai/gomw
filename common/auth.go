package common

import (
	"log"
	"net/http"

	"github.com/vchitai/gomw"
)

type AuthenticateFunc func(request *http.Request) (*http.Request, error)

func AuthenticationHTTPMiddleware(authenticateFunc AuthenticateFunc) gomw.HTTPMiddleware {
	return gomw.NewHTTPBeforeMiddleware(func(writer http.ResponseWriter, request *http.Request) (*http.Request, bool) {
		request, err := authenticateFunc(request)
		if err != nil {
			log.Printf("Validation error: %v\n", err)
			// terminate the request instantly
			writer.WriteHeader(http.StatusUnauthorized)
			_, _ = writer.Write([]byte(http.StatusText(http.StatusUnauthorized)))
			return request, false
		}

		return request, true
	})
}
