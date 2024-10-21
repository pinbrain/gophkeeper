package interceptors

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pinbrain/gophkeeper/internal/client/config"
	pb "github.com/pinbrain/gophkeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TokenInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		jwt := config.GetJWT()
		if jwt != "" && strings.Split(method, "/")[1] != pb.UserService_ServiceDesc.ServiceName {
			md := metadata.Pairs(config.GetJWTMetaKey(), jwt)
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			if s, ok := status.FromError(err); ok && s.Code() == codes.Unauthenticated {
				fmt.Println("Unauthenticated error received. Deleting token.")

				if jwtErr := config.SaveJWT(""); jwtErr != nil {
					log.Fatalf("Error clearing jwt from config file: %v", jwtErr)
				}
			}
		}

		return err
	}
}
