//go:build !solution

package requestlog

import (
	"net/http"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

func Log(log *zap.Logger) func(nextHandler http.Handler) http.Handler {
	var errOccurred any
	var request *http.Request
	var duration time.Duration
	var requestID uuid.UUID
	var responseStatus int

	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			id, _ := uuid.NewV4()
			requestID = id
			request = req

			log.Info("request started",
				zap.String("method", request.Method),
				zap.String("path", request.RequestURI),
				zap.String("request_id", requestID.String()),
			)

			writer = httpsnoop.Wrap(writer,
				httpsnoop.Hooks{WriteHeader: func(originalHeader httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
					return func(status int) {
						responseStatus = status
						originalHeader(status)
					}
				}})

			func() {
				defer func() {
					errOccurred = recover()
				}()

				startTime := time.Now()
				nextHandler.ServeHTTP(writer, req)
				duration = time.Since(startTime)
			}()

			if errOccurred != nil {
				log.Info("request panicked",
					zap.String("path", request.RequestURI),
					zap.String("method", request.Method),
					zap.Float64("duration", duration.Seconds()),
					zap.Int("status_code", responseStatus),
					zap.String("request_id", requestID.String()),
				)
				panic(errOccurred)
			} else {
				log.Info("request finished",
					zap.String("path", request.RequestURI),
					zap.String("method", request.Method),
					zap.Int("status_code", responseStatus),
					zap.Float64("duration", duration.Seconds()),
					zap.String("request_id", requestID.String()),
				)
			}
		})
	}
}
