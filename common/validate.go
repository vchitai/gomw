package common

import (
	"log"
	"net/http"

	"github.com/vchitai/gomw"
)

type ValidateFunc func(request *http.Request) error

func ValidationHTTPMiddleware(validateFunc ValidateFunc) gomw.HTTPMiddleware {
	return gomw.NewHTTPBeforeMiddleware(func(writer http.ResponseWriter, request *http.Request) (*http.Request, bool) {
		if err := validateFunc(request); err != nil {
			log.Printf("Validation error: %v\n", err)
			// terminate the request instantly
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return request, false
		}

		return request, true
	})
}
