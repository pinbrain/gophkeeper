package interceptors

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// LoggerInterceptor логирует входящие запросы.
func LoggerInterceptor(
	log *logrus.Entry,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(start).Seconds()
		status, _ := status.FromError(err)

		log.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"duration": duration,
			"code":     status.Code().String(),
		}).Info("gRPC request")
		return resp, err
	}
}
