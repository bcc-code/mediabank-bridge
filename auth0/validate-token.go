package auth0

import (
	"github.com/lestrrat-go/jwx/jwt"
)

// JWTConfig for configuring what passes as a valid JWT
type JWTConfig struct {
	Domain   string
	Issuer   string
	Audience string
}

// ValidateJWT and return an error if invalid
// No values are extracted in this process
func ValidateJWT(token string, config JWTConfig) error {
	jwks := GetKeySet(config.Domain)

	t, err := jwt.ParseString(token,
		jwt.WithKeySet(jwks),
	)
	if err != nil {
		return err
	}

	return jwt.Validate(
		t,
		jwt.WithIssuer(config.Issuer),
		jwt.WithAudience(config.Audience),
	)
}
