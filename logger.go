package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type contextKey int

const loggerKey contextKey = iota

func init() {
	defaultLogger = logrus.New()
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corrID := uuid.New()
		logger := defaultLogger.WithField("correlationid", corrID)

		ctx := r.Context()
		ctx = context.WithValue(ctx, loggerKey, logger)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getLogger(ctx context.Context) *logrus.Entry {
	logger := ctx.Value(loggerKey)
	if logger == nil {
		logrus.Fatal("Logger missing")
	}

	tlogger, ok := logger.(*logrus.Entry)
	if !ok {
		logrus.Fatal("Not a *logrus.Logger")
	}

	return tlogger
}
