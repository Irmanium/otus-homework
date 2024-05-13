package service

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *Service) getJWTConfig() echojwt.Config {
	return echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte(s.jwtSecret),
	}
}

func (s *Service) getUserIDFromJWT(c echo.Context) string {
	return c.Get("user").(*jwt.Token).Claims.(*jwtCustomClaims).UserID
}

func (s *Service) generateToken(id string) (string, error) {
	claims := &jwtCustomClaims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(s.tokenTTLHours))),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return t, nil
}
