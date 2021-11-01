package common

import (
	"log"
	"net/http"

	"github.com/vchitai/gomw"
)

func PayloadLoggingHTTPMiddleware() gomw.HTTPMiddleware {
	return gomw.NewHTTPMiddleware(func(writer http.ResponseWriter, request *http.Request) (*http.Request, bool) {
		bodyReader, err := request.GetBody()
		if err != nil {
			// not logging, pass through this step
			return request, true
		}
		var body []byte
		if _, err := bodyReader.Read(body); err != nil {
			// not logging, pass through this step
			return request, true
		}
		log.Println("A request was recorded", "url", request.URL.String(), "payload", string(body))
		return request, true
	}, func(response gomw.HTTPResponse, request *http.Request) gomw.HTTPResponse {
		log.Println("A response was recorded", "url", request.URL.String(), "payload", string(response.Body()), "code", response.Code())
		return response
	})
}
