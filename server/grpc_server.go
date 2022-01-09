package server

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/taehoio/auth/config"
	authv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1"
)

type AuthServiceServer struct {
	authv1.AuthServiceServer

	cfg config.Config
}

func NewAuthServiceServer(cfg config.Config) (*AuthServiceServer, error) {
	return &AuthServiceServer{
		cfg: cfg,
	}, nil
}

func (s *AuthServiceServer) HealthCheck(ctx context.Context, req *authv1.HealthCheckRequest) (*authv1.HealthCheckResponse, error) {
	return &authv1.HealthCheckResponse{}, nil
}

func (s *AuthServiceServer) AuthByRefreshToken(ctx context.Context, req *authv1.AuthByRefreshTokenRequest) (*authv1.AuthByRefreshTokenResponse, error) {
	return &authv1.AuthByRefreshTokenResponse{}, nil
}

func (s *AuthServiceServer) VerifyToken(ctx context.Context, req *authv1.VerifyTokenRequest) (*authv1.VerifyTokenResponse, error) {
	return &authv1.VerifyTokenResponse{}, nil
}

func (s *AuthServiceServer) ParseToken(ctx context.Context, req *authv1.ParseTokenRequest) (*authv1.ParseTokenResponse, error) {
	return &authv1.ParseTokenResponse{}, nil
}

func (s *AuthServiceServer) Auth(ctx context.Context, req *authv1.AuthRequest) (*authv1.AuthResponse, error) {
	return &authv1.AuthResponse{}, nil
}

func NewGRPCServer(cfg config.Config) (*grpc.Server, error) {
	logrus.ErrorKey = "grpc.error"
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(
				grpc_ctxtags.WithFieldExtractor(
					grpc_ctxtags.CodeGenRequestFieldExtractor,
				),
			),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 30 * time.Second,
		}),
	)

	authServiceServer, err := NewAuthServiceServer(cfg)
	if err != nil {
		return nil, err
	}

	authv1.RegisterAuthServiceServer(grpcServer, authServiceServer)
	reflection.Register(grpcServer)

	return grpcServer, nil
}