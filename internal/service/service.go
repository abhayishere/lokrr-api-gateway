package service

import (
	"github.com/abhayishere/lokrr-proto/gen/authpb"
	fmpb "github.com/abhayishere/lokrr-proto/gen/file_management"
	"google.golang.org/grpc"
)

type serviceImpl struct {
	AuthClient authpb.AuthServiceClient
	DocClient  fmpb.FileManagementServiceClient
}

func NewService(authConn, docConn *grpc.ClientConn) *serviceImpl {
	return &serviceImpl{
		authpb.NewAuthServiceClient(authConn),
		fmpb.NewFileManagementServiceClient(docConn),
	}
}
