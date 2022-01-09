package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testHMACSecret = "abcd1234"
	testIssuer     = "taeho.io"
	testAudience   = "taeho.io"
)

func TestSignAndParse(t *testing.T) {
	j := NewHS256JWT(testHMACSecret, testIssuer, testAudience)

	claims := map[string]interface{}{
		"provider":   "email",
		"identifier": "taeho@taeho.io",
		"token_type": "access_token",
	}
	s, err := j.Sign(time.Minute*15, claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	c, err := j.Parse(s)
	assert.NoError(t, err)
	assert.Equal(t, testAudience, c["aud"])
	assert.Equal(t, testIssuer, c["iss"])
	assert.Equal(t, "email", c["provider"])
	assert.Equal(t, "taeho@taeho.io", c["identifier"])
	assert.Equal(t, "access_token", c["token_type"])
}

func TestParse_FailedWithExpiredToken(t *testing.T) {
	j := NewHS256JWT(testHMACSecret, testIssuer, testAudience)

	claims := map[string]interface{}{}
	s, err := j.Sign(0, claims)
	assert.NoError(t, err)
	assert.NotEmpty(t, s)

	c, err := j.Parse(s)
	assert.Error(t, err)
	assert.Nil(t, c)
}
