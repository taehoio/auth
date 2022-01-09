package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWT struct {
	signingMethod *jwt.SigningMethodHMAC
	hmacSecret    string
	issuer        string
	audience      string
}

func NewHS256JWT(hmacSecret, issuer, audience string) *JWT {
	return &JWT{
		signingMethod: jwt.SigningMethodHS256,
		hmacSecret:    hmacSecret,
		issuer:        issuer,
		audience:      audience,
	}
}

func (j *JWT) Sign(expiresIn time.Duration, claims map[string]interface{}) (string, error) {
	currentAt := time.Now()

	jc := make(jwt.MapClaims)
	jc["iss"] = j.issuer
	jc["aud"] = j.audience
	jc["iat"] = currentAt.Unix()
	jc["exp"] = currentAt.Add(expiresIn).Unix()

	for k, v := range claims {
		jc[k] = v
	}

	token := jwt.NewWithClaims(j.signingMethod, jc)
	return token.SignedString([]byte(j.hmacSecret))
}

func (j *JWT) Parse(s string) (map[string]interface{}, error) {
	tkn, err := jwt.Parse(s, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(j.hmacSecret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
