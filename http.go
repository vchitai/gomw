package gomw

import (
	"net/http"

	"go.uber.org/zap/buffer"
)

type HTTPMiddleware func(next http.Handler) http.Handler

var _ http.ResponseWriter = &copyWriter{}
var _ HTTPResponse = &copyWriter{}
var bPool = buffer.NewPool()

type copyWriter struct {
	target http.ResponseWriter
	buf    *buffer.Buffer
	code   int
}

func newCopyWriter(target http.ResponseWriter) *copyWriter {
	return &copyWriter{
		target: target,
		buf:    bPool.Get(),
	}
}

func (m *copyWriter) Body() []byte {
	return m.buf.Bytes()
}

func (m *copyWriter) Header() http.Header {
	return m.target.Header()
}

func (m *copyWriter) Code() int {
	if m.code == 0 {
		return http.StatusOK
	}
	return m.code
}

func (m *copyWriter) Write(i []byte) (int, error) {
	return m.buf.Write(i)
}

func (m *copyWriter) WriteHeader(statusCode int) {
	m.code = statusCode
}

// push the buffer to writer
func (m *copyWriter) free() {
	m.buf.Free()
}

type httpResponse struct {
	body []byte
	code int
}

func (h *httpResponse) Body() []byte {
	return h.body
}

func (h *httpResponse) Code() int {
	if h.code == 0 {
		return http.StatusOK
	}
	return h.code
}

func NewHTTPResponse(body []byte, code int) HTTPResponse {
	return &httpResponse{
		body: body,
		code: code,
	}
}

func afterMiddleware(after AfterHook) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			var copyWriter = newCopyWriter(writer)
			defer copyWriter.free()
			next.ServeHTTP(copyWriter, request)
			var resp = after(copyWriter, request)
			writer.WriteHeader(resp.Code())
			_, _ = writer.Write(resp.Body())
		})
	}
}

func beforeMiddleware(before BeforeHook) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request, ok := before(writer, request)
			if !ok {
				return
			}
			next.ServeHTTP(writer, request)
		})
	}
}

func fullyMiddleware(before BeforeHook, after AfterHook) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request, ok := before(writer, request)
			if !ok {
				return
			}
			var copyWriter = newCopyWriter(writer)
			defer copyWriter.free()
			next.ServeHTTP(copyWriter, request)
			var resp = after(copyWriter, request)
			writer.WriteHeader(resp.Code())
			_, _ = writer.Write(resp.Body())
		})
	}
}

func doNothingMiddleware() HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return next // nothing to do
	}
}

type HTTPResponse interface {
	Body() []byte
	Code() int
}
type BeforeHook func(writer http.ResponseWriter, request *http.Request) (*http.Request, bool)
type AfterHook func(response HTTPResponse, request *http.Request) HTTPResponse

func NewHTTPMiddleware(before BeforeHook, after AfterHook) HTTPMiddleware {
	if before == nil && after == nil {
		return doNothingMiddleware()
	}

	if before == nil {
		return afterMiddleware(after)
	}

	if after == nil {
		return beforeMiddleware(before)
	}

	return fullyMiddleware(before, after)
}

func NewHTTPBeforeMiddleware(before BeforeHook) HTTPMiddleware {
	if before == nil {
		return doNothingMiddleware()
	}
	return beforeMiddleware(before)
}

func NewHTTPAfterMiddleware(after AfterHook) HTTPMiddleware {
	if after == nil {
		return doNothingMiddleware()
	}
	return afterMiddleware(after)
}
