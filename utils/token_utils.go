package utils

import (
	"errors"
	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"interviewexcel-backend-go/config"
	logger "interviewexcel-backend-go/pkg/errors"
	"os"
	"time"
)

// ❗ Move secrets to ENV (config.JWTAccessSecret, config.JWTRefreshSecret)
var (
	accessSecret  = []byte(os.Getenv("JWT_SECRET")) // shorter lifetime
	refreshSecret = []byte(os.Getenv("JWT_SECRET")) // longer lifetime
)

type Claims struct {
	UserID string `json:"user_uuid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func getAccessSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logger.Fatal("JWT_SECRET not set")
	}
	logger.Info("secret",secret)
	return []byte(secret)
}

func GenerateAccessToken(userID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "access_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getAccessSecret())
}

// GenerateRefreshToken issues a refresh token (default 30 days)
func GenerateRefreshToken(userID string, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "refresh_token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getAccessSecret())
}

func ValidateAccessToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return accessSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}
	return claims, nil
}

func ValidateRefreshToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}
	return claims, nil
}

// AddTokenToBlacklist adds a token to Redis blacklist with expiration
func AddTokenToBlacklist(token string, expiration time.Duration) error {
	return config.RedisClient.Set(config.Ctx, token, true, expiration).Err()
}

// IsTokenBlacklisted checks if a token exists in Redis blacklist
func IsTokenBlacklisted(token string) (bool, error) {
	exists, err := config.RedisClient.Exists(config.Ctx, token).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

var RefreshTokenSecret = []byte(os.Getenv("JWT_SECRET")) // use env variable in production

type RefreshClaims struct {
	UserID string `json:"user_uuid"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// VerifyRefreshToken validates the refresh token and returns the claims if valid
func VerifyRefreshToken(tokenStr string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return RefreshTokenSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	// Optional: check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("refresh token expired")
	}

	return claims, nil
}
