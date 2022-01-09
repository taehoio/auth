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
)

var (
	ErrInvalidRefreshToken = status.Error(codes.InvalidArgument, "invalid refresh_token")
	ErrInvalidTokenType    = status.Error(codes.InvalidArgument, "invalid token_type")
)

type AuthByRefreshTokenHandlerFunc func(ctx context.Context, req *authv1.AuthByRefreshTokenRequest) (*authv1.AuthByRefreshTokenResponse, error)

func AuthByRefreshTokenHandler(cfg config.Config, jwt *jwt.JWT) AuthByRefreshTokenHandlerFunc {
	return func(ctx context.Context, req *authv1.AuthByRefreshTokenRequest) (*authv1.AuthByRefreshTokenResponse, error) {
		if err := validateAuthByRefreshTokenRequest(req); err != nil {
			return nil, err
		}

		provider, identifier, tokenType, err := parseToken(jwt, req.RefreshToken)
		if err != nil {
			return nil, err
		}

		if tokenType != authv1.TokenType_TOKEN_TYPE_REFRESH {
			return nil, ErrInvalidTokenType
		}

		accessTokenExpiresIn := time.Duration(cfg.Setting().AccessTokenExpiresInMs) * time.Millisecond
		accessToken, err := newToken(
			jwt,
			accessTokenExpiresIn,
			provider,
			identifier,
			authv1.TokenType_TOKEN_TYPE_ACCESS,
		)
		if err != nil {
			return nil, err
		}

		refreshTokenExpiresIn := time.Duration(cfg.Setting().RefreshTokenExpiresInMs) * time.Millisecond
		refreshToken, err := newToken(
			jwt,
			accessTokenExpiresIn,
			provider,
			identifier,
			authv1.TokenType_TOKEN_TYPE_REFRESH,
		)
		if err != nil {
			return nil, err
		}

		return &authv1.AuthByRefreshTokenResponse{
			AccessToken:           accessToken,
			AccessTokenExpiresIn:  durationpb.New(accessTokenExpiresIn),
			RefreshToken:          refreshToken,
			RefreshTokenExpiresIn: durationpb.New(refreshTokenExpiresIn),
		}, nil
	}
}

func validateAuthByRefreshTokenRequest(req *authv1.AuthByRefreshTokenRequest) error {
	if req.RefreshToken == "" {
		return ErrInvalidRefreshToken
	}
	return nil
}
