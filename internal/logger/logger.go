package logger

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"time"

	ip "github.com/vikram1565/request-ip"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const ErrorField = "error"

type loggerCtxKey struct{}

func DefaultLogger() *zap.Logger {
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig()), os.Stdout, zap.DebugLevel))
}
func WithLogger(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey{}, log)
}
func GetLogger(ctx context.Context) *zap.Logger {
	le, ok := ctx.Value(loggerCtxKey{}).(*zap.Logger)
	if !ok {
		le = DefaultLogger()
	}
	return le
}

type statusWriter struct {
	http.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func NewLoggingMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			methodLog := zapcore.Field{
				Key:    "method",
				Type:   zapcore.StringType,
				String: r.Method,
			}
			pathLog := zapcore.Field{
				Key:    "path",
				Type:   zapcore.StringType,
				String: r.URL.Path,
			}
			ipLog := zapcore.Field{
				Key:    "ip",
				Type:   zapcore.StringType,
				String: ip.GetClientIP(r),
			}
			nextLog := log.With(
				methodLog,
				pathLog,
				ipLog,
			)

			log.Info("Request received", methodLog, pathLog, ipLog)
			sw := &statusWriter{ResponseWriter: w}

			startTime := time.Now()
			handler.ServeHTTP(sw, r.WithContext(WithLogger(r.Context(), nextLog)))
			elapseTime := time.Since(startTime)

			statusLog := zapcore.Field{
				Key:     "status",
				Type:    zapcore.Int64Type,
				Integer: int64(sw.status),
			}
			timeLog := zapcore.Field{
				Key:     "time",
				Type:    zapcore.DurationType,
				Integer: int64(elapseTime),
			}

			if sw.status >= 200 && sw.status < 300 {
				nextLog.Info("Request success", timeLog, statusLog)
			} else {
				nextLog.Info("Request failed", timeLog, statusLog, zapcore.Field{
					Key:    "error response",
					Type:   zapcore.StringType,
					String: sw.body.String(),
				})
			}
		})
	}
}
