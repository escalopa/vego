package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/escalopa/vego/internal/config"
	"github.com/escalopa/vego/internal/domain"
)

var getCurrentTime = time.Now

type (
	RoomProvider struct {
		secretKey []byte
		tokenTTL  time.Duration
	}

	roomClaims struct {
		UserID int64  `json:"user_id"`
		RoomID string `json:"room_id"`
		jwt.RegisteredClaims
	}
)

func NewRoomProvider(cfg config.JWTRoom) *RoomProvider {
	return &RoomProvider{
		secretKey: []byte(cfg.SecretKey),
		tokenTTL:  cfg.TokenTTL,
	}
}

func (rp *RoomProvider) CreateToken(userID int64, roomID string) (string, error) {
	now := time.Now()
	claims := roomClaims{
		UserID: userID,
		RoomID: roomID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(rp.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(rp.secretKey)
}

func (rp *RoomProvider) VerifyToken(tokenStr string) (*domain.RoomTokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &roomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return rp.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired
		}
		return nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*roomClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	payload := &domain.RoomTokenPayload{
		UserID: claims.UserID,
		RoomID: claims.RoomID,
	}
	return payload, nil
}

type (
	UserProvider struct {
		secretKey       []byte
		accessTokenTTL  time.Duration
		refreshTokenTTL time.Duration
	}

	userClaims struct {
		UserID int64  `json:"user_id"`
		Email  string `json:"email"`
		jwt.RegisteredClaims
	}
)

func NewUserProvider(cfg config.JWTUser) *UserProvider {
	return &UserProvider{
		secretKey:       []byte(cfg.SecretKey),
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
	}
}

func (up *UserProvider) CreateToken(userID int64, email string) (*domain.Token, error) {
	now := time.Now()
	accessToken, err := up.createToken(userID, email, up.accessTokenTTL, now)
	if err != nil {
		return nil, err
	}

	refreshToken, err := up.createToken(userID, email, up.refreshTokenTTL, now)
	if err != nil {
		return nil, err
	}

	token := &domain.Token{
		Access:  accessToken,
		Refresh: refreshToken,
	}

	return token, nil
}

func (up *UserProvider) createToken(userID int64, email string, ttl time.Duration, now time.Time) (string, error) {
	claims := userClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(up.secretKey)
}

func (up *UserProvider) VerifyToken(tokenStr string) (*domain.UserTokenPayload, error) {
	if tokenStr == "" {
		return nil, domain.ErrTokenExpired // treat empty token as expired
	}

	token, err := jwt.ParseWithClaims(tokenStr, &userClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return up.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrTokenExpired
		}
		return nil, domain.ErrTokenInvalid
	}

	claims, ok := token.Claims.(*userClaims)
	if !ok || !token.Valid {
		return nil, domain.ErrTokenInvalid
	}

	payload := &domain.UserTokenPayload{
		UserID: claims.UserID,
		Email:  claims.Email,
	}
	return payload, nil
}
