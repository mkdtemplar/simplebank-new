package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Write(body []byte) (int, error) {
	r.Body = body
	return r.ResponseWriter.Write(body)
}

func GrpcLogger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "gRPC").Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).Str("status_text", statusCode.String()).
		Dur("duration", duration).Msg("gRPC Logger: ")
	return result, err
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{ResponseWriter: w, StatusCode: http.StatusOK}
		logger := log.Info()

		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		handler.ServeHTTP(rec, r)
		duration := time.Since(startTime)

		logger.Str("protocol", "http").Str("method", r.Method).Str("path", r.RequestURI).
			Int("status_code", rec.StatusCode).Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).Msg("HTTP Logger: ")
	})
}
