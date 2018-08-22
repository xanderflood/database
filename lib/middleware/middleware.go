package middleware

import (
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/xanderflood/database/lib/tools"
)

//Wrap wrap with standard middleware chain
func Wrap(log tools.Logger, h http.Handler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error("500 response failed: forgoing 500")
					log.Errorf("Panic: %v", r)
					spew.Dump(r)
				}
			}()

			if r := recover(); r != nil {
				log.Error("Request failed: 500")
				log.Errorf("Panic: %v", r)
				spew.Dump(r)

				http.Error(w, http.StatusText(500), 500)
			}
		}()

		log.Infof("Request: %s %s", r.Method, r.URL.String())

		lrw := &LoggingResponseWriter{w: w, log: log}
		h.ServeHTTP(lrw, r)
	}
}

//////

//LoggingResponseWriter wraps http.ResponseWriter with a logger
type LoggingResponseWriter struct {
	w   http.ResponseWriter
	log tools.Logger
}

//Header Header
func (lrw *LoggingResponseWriter) Header() http.Header {
	return lrw.w.Header()
}

//Write Write
func (lrw *LoggingResponseWriter) Write(p []byte) (int, error) {
	return lrw.w.Write(p)
}

//WriteHeader WriteHeader
func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.log.Infof("Response: %v", statusCode)

	lrw.w.WriteHeader(statusCode)
}
