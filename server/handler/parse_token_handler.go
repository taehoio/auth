package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/taehoio/auth/config"
	"github.com/taehoio/auth/internal/jwt"
	authv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1"
	userv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/user/v1"
)

var (
	ErrInvalidToken = status.Error(codes.InvalidArgument, "invalid token")
)

type ParseTokenHandlerFunction func(ctx context.Context, req *authv1.ParseTokenRequest) (*authv1.ParseTokenResponse, error)

func ParseToken(cfg config.Config, jwt *jwt.JWT) ParseTokenHandlerFunction {
	return func(ctx context.Context, req *authv1.ParseTokenRequest) (*authv1.ParseTokenResponse, error) {
		if validateParseTokenRequest(req) != nil {
			return nil, ErrInvalidToken
		}

		provider, identifier, tokenType, err := parseToken(jwt, req.Token)
		if err != nil {
			return nil, err
		}

		return &authv1.ParseTokenResponse{
			Provider:   provider,
			Identifier: identifier,
			TokenType:  tokenType,
		}, nil
	}
}

func validateParseTokenRequest(req *authv1.ParseTokenRequest) error {
	if req.Token == "" {
		return ErrInvalidToken
	}

	return nil
}

func parseToken(jwt *jwt.JWT, s string) (userv1.Provider, string, authv1.TokenType, error) {
	claims, err := jwt.Parse(s)
	if err != nil {
		return userv1.Provider_PROVIDER_UNSPECIFIED, "", authv1.TokenType_TOKEN_TYPE_UNSPECIFIED, err
	}

	provider := userv1.Provider(claims["provider"].(float64))
	identifier := claims["identifier"].(string)
	tokenType := authv1.TokenType(claims["token_type"].(float64))

	return provider, identifier, tokenType, nil
}
