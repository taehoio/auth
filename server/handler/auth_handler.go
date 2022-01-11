package handler

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/taehoio/auth/config"
	"github.com/taehoio/auth/internal/jwt"
	authv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1"
	userv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/user/v1"
)

var (
	ErrInvalidProvider   = status.Error(codes.InvalidArgument, "invalid provider")
	ErrInvalidIdentifier = status.Error(codes.InvalidArgument, "invalid identifier")
)

type AuthHandlerFunc func(ctx context.Context, req *authv1.AuthRequest) (*authv1.AuthResponse, error)

func Auth(cfg config.Config, jwt *jwt.JWT) AuthHandlerFunc {
	return func(ctx context.Context, req *authv1.AuthRequest) (*authv1.AuthResponse, error) {

		if err := validateAuthRequest(req); err != nil {
			return nil, err
		}

		accessTokenExpiresIn := time.Duration(cfg.Setting().AccessTokenExpiresInMs) * time.Millisecond
		accessToken, err := newToken(
			jwt,
			accessTokenExpiresIn,
			req.Provider,
			req.Identifier,
			authv1.TokenType_TOKEN_TYPE_ACCESS,
		)
		if err != nil {
			return nil, err
		}

		refreshTokenExpiresIn := time.Duration(cfg.Setting().RefreshTokenExpiresInMs) * time.Millisecond
		refreshToken, err := newToken(
			jwt,
			accessTokenExpiresIn,
			req.Provider,
			req.Identifier,
			authv1.TokenType_TOKEN_TYPE_REFRESH,
		)
		if err != nil {
			return nil, err
		}

		return &authv1.AuthResponse{
			AccessToken:           accessToken,
			AccessTokenExpiresIn:  durationpb.New(accessTokenExpiresIn),
			RefreshToken:          refreshToken,
			RefreshTokenExpiresIn: durationpb.New(refreshTokenExpiresIn),
		}, nil
	}
}

func validateAuthRequest(req *authv1.AuthRequest) error {
	if req.Provider == userv1.Provider_PROVIDER_UNSPECIFIED {
		return ErrInvalidProvider
	}
	if req.Identifier == "" {
		return ErrInvalidIdentifier
	}
	return nil
}

func newToken(
	jwt *jwt.JWT,
	expiresIn time.Duration,
	provider userv1.Provider,
	identifier string,
	tokenType authv1.TokenType,
) (string, error) {
	claims := map[string]interface{}{
		"provider":   provider,
		"identifier": identifier,
		"token_type": tokenType,
	}
	return jwt.Sign(expiresIn, claims)
}
