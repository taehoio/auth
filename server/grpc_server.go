package server

import (
	"context"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/taehoio/auth/config"
	"github.com/taehoio/auth/internal/jwt"
	"github.com/taehoio/auth/server/handler"
	authv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1"
)

type AuthServiceServer struct {
	authv1.AuthServiceServer

	cfg config.Config
	jwt *jwt.JWT
}

func NewAuthServiceServer(cfg config.Config) (*AuthServiceServer, error) {
	return &AuthServiceServer{
		cfg: cfg,
		jwt: jwt.NewHS256JWT(
			cfg.Setting().JWTHMACSecret,
			cfg.Setting().JWTIssuer,
			cfg.Setting().JWTAudience,
		),
	}, nil
}

func (s *AuthServiceServer) HealthCheck(ctx context.Context, req *authv1.HealthCheckRequest) (*authv1.HealthCheckResponse, error) {
	return &authv1.HealthCheckResponse{}, nil
}

func (s *AuthServiceServer) AuthByRefreshToken(ctx context.Context, req *authv1.AuthByRefreshTokenRequest) (*authv1.AuthByRefreshTokenResponse, error) {
	return handler.AuthByRefreshTokenHandler(s.cfg, s.jwt)(ctx, req)
}

func (s *AuthServiceServer) Auth(ctx context.Context, req *authv1.AuthRequest) (*authv1.AuthResponse, error) {
	return handler.Auth(s.cfg, s.jwt)(ctx, req)
}

func NewGRPCServer(cfg config.Config) (*grpc.Server, error) {
	logrus.ErrorKey = "grpc.error"
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())

	grpcServer := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
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

func (s *AuthServiceServer) ParseToken(ctx context.Context, req *authv1.ParseTokenRequest) (*authv1.ParseTokenResponse, error) {
	return &authv1.ParseTokenResponse{}, nil
}
