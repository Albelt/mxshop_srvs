package middleware

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const (
	grpcHealthCheckMethodName = "/grpc.health.v1.Health/Check"
)

func ServerLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startT := time.Now()
		resp, err = handler(ctx, req)
		endT := time.Now()

		if info.FullMethod != grpcHealthCheckMethodName {
			// 打印grpc方法和请求参数
			bytes, _ := json.Marshal(req)
			zap.S().Infof("[grpc] method: %s, latency:%s, request: %s, ",
				info.FullMethod, endT.Sub(startT).String(), string(bytes))

			// 打印错误信息
			if err != nil {
				var (
					code, msg string
				)

				if s, ok := status.FromError(err); ok {
					code, msg = s.Code().String(), s.Message()
				} else {
					code, msg = codes.Unknown.String(), err.Error()
				}

				zap.S().Infof("[grpc] error occured, code: %s, msg: %s", code, msg)
			}
		}

		return resp, err
	}
}
