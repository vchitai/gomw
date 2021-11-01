package common

import (
	"net/http"

	"github.com/vchitai/gomw"
)

func ErrorHidingHTTPMiddleware() gomw.HTTPMiddleware {
	return gomw.NewHTTPAfterMiddleware(func(response gomw.HTTPResponse, request *http.Request) gomw.HTTPResponse {
		if response.Code() == http.StatusInternalServerError {
			return gomw.NewHTTPResponse([]byte(http.StatusText(http.StatusInternalServerError)), http.StatusInternalServerError)
		}
		return response
	})
}
