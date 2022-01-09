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

const (
	testHMACSecret = "testHMACSecret"
	testIssuer     = "testIssuer"
	testAudience   = "testAudience"
)

func TestAuth(t *testing.T) {
	cfg := config.MockConfig()
	jwt := jwt.NewHS256JWT(testHMACSecret, testIssuer, testAudience)

	req := &authv1.AuthRequest{
		Provider:   userv1.Provider_PROVIDER_EMAIL,
		Identifier: "taeho@taeho.io",
	}
	resp, err := Auth(cfg, jwt)(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}
