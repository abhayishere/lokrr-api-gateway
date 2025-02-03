package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/abhayishere/lokrr-api-gateway/internal/client"
	"github.com/abhayishere/lokrr-api-gateway/internal/config"
	"github.com/abhayishere/lokrr-api-gateway/internal/middlerware"
	"github.com/abhayishere/lokrr-api-gateway/pkg"
	proto "github.com/abhayishere/lokrr-auth-service/proto/authService"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

func main() {
	cfg, err := config.LoadConfig()
	fmt.Println(cfg.Grpc.AuthServiceAddress)
	if err != nil {
		pkg.LogError(err)
		panic(err)
	}

	mux := runtime.NewServeMux()
	authConn, err := client.NewAuthServiceClient(cfg.Grpc.AuthServiceAddress)
	if err != nil {
		pkg.LogError(err)
		panic(err)
	}
	err = proto.RegisterAuthServiceHandler(context.Background(), mux, authConn)
	if err != nil {
		pkg.LogError(err)
		panic(err)
	}
	pkg.LogInfo("API Gateway running on http://localhost:8080...")
	wrappedMux := middlerware.LoggingMiddleware(mux)
	http.ListenAndServe(":8080", wrappedMux)
}
