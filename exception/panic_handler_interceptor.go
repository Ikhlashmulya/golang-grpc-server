package exception

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func PanicHandlerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		defer func() {
			if recover := recover(); recover != nil {
				log.Printf("Recovered from panic: %v", recover)
				grpc.SetHeader(ctx, metadata.Pairs("panic", "true"))
				grpc.SendHeader(ctx, metadata.Pairs("panic", "true"))
				grpc.SetTrailer(ctx, metadata.Pairs("panic", "true"))

				if errors.Is(sql.ErrNoRows, recover.(error)) {
					err = status.Errorf(codes.NotFound, sql.ErrNoRows.Error())
					return
				}

				err = status.Errorf(codes.Internal, "Internal server error")
				return
			}
		}()
		return handler(ctx, req)
	}
}
