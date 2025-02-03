package client

import (
	"github.com/abhayishere/lokrr-api-gateway/pkg"
	"google.golang.org/grpc"
)

func NewAuthServiceClient(grpcAuthServer string) (*grpc.ClientConn, error) {
	authConn, err := grpc.NewClient(grpcAuthServer, grpc.WithInsecure())
	if err != nil {
		pkg.LogError(err)
		panic(err)
	}
	return authConn, nil
}
