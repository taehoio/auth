package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/taehoio/auth/config"
	"github.com/taehoio/auth/internal/jwt"
	authv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/auth/v1"
	userv1 "github.com/taehoio/idl/gen/go/taehoio/idl/services/user/v1"
)

func TestParseToken(t *testing.T) {
	cfg := config.MockConfig()
	jwt := jwt.NewHS256JWT(testHMACSecret, testIssuer, testAudience)

	authReq := &authv1.AuthRequest{
		Provider:   userv1.Provider_PROVIDER_EMAIL,
		Identifier: "taeho@taeho.io",
	}
	authResp, err := Auth(cfg, jwt)(context.Background(), authReq)
	assert.NoError(t, err)

	t.Run("withAccessToken", func(t *testing.T) {
		req := &authv1.ParseTokenRequest{
			Token: authResp.AccessToken,
		}
		resp, err := ParseToken(cfg, jwt)(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, userv1.Provider_PROVIDER_EMAIL, resp.Provider)
		assert.Equal(t, "taeho@taeho.io", resp.Identifier)
		assert.Equal(t, authv1.TokenType_TOKEN_TYPE_ACCESS, resp.TokenType)
	})

	t.Run("withRefreshToken", func(t *testing.T) {
		req := &authv1.ParseTokenRequest{
			Token: authResp.RefreshToken,
		}
		resp, err := ParseToken(cfg, jwt)(context.Background(), req)
		assert.NoError(t, err)
		assert.Equal(t, userv1.Provider_PROVIDER_EMAIL, resp.Provider)
		assert.Equal(t, "taeho@taeho.io", resp.Identifier)
		assert.Equal(t, authv1.TokenType_TOKEN_TYPE_REFRESH, resp.TokenType)
	})
}
