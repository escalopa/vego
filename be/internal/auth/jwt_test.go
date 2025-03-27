package auth

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/escalopa/vego/internal/config"
	"github.com/escalopa/vego/internal/domain"
	"github.com/stretchr/testify/require"
)

const (
	testUserID = int64(1)
	testEmail  = "test@example.com"
	testRoomID = "room1"
)

var (
	expiredToken, _ = base64.RawStdEncoding.DecodeString("ZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y" +
		"0NJNklrcFhWQ0o5LmV5SjFjMlZ5WDJsa0lqb3hMQ0p5YjI5dFgybGtJam9pY205dmJURW" +
		"lMQ0psZUhBaU9qRTNNVEUwT1RJNU1EUXNJbWxoZENJNk1UY3hNVFE1TWpnME5IMC5YYmh" +
		"hR2g4VUNTdmp5M1UxRkVEV3dmVl9qYmU1dkVNU29oVEVtemswMDVn",
	)
)

func TestRoomProvider(t *testing.T) {
	t.Parallel()

	cfg := config.JWTRoom{
		SecretKey: "test-secret",
		TokenTTL:  time.Minute,
	}
	p := NewRoomProvider(cfg)

	tests := []struct {
		name      string
		userID    int64
		roomID    string
		modify    func(token string) string
		expectErr error
	}{
		{
			name:      "valid_token",
			userID:    testUserID,
			roomID:    testRoomID,
			modify:    func(token string) string { return token },
			expectErr: nil,
		},
		{
			name:      "expired_token",
			userID:    testUserID,
			roomID:    testRoomID,
			modify:    func(_ string) string { return string(expiredToken) },
			expectErr: domain.ErrTokenExpired,
		},
		{
			name:      "invalid_token",
			userID:    testUserID,
			roomID:    testRoomID,
			modify:    func(token string) string { return token + "invalid" },
			expectErr: domain.ErrTokenInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token, err := p.CreateToken(tt.userID, tt.roomID)
			require.NoError(t, err)

			token = tt.modify(token)
			payload, err := p.VerifyToken(token)

			if tt.expectErr == nil {
				require.NoError(t, err)
				require.NotNil(t, payload)
				require.Equal(t, tt.userID, payload.UserID)
				require.Equal(t, tt.roomID, payload.RoomID)
			} else {
				require.ErrorIs(t, err, tt.expectErr)
				require.Nil(t, payload)
			}
		})
	}
}

func TestAuthProvider(t *testing.T) {
	t.Parallel()

	cfg := config.JWTUser{
		SecretKey:       "test-secret",
		AccessTokenTTL:  time.Minute,
		RefreshTokenTTL: time.Hour,
	}
	p := NewUserProvider(cfg)

	tests := []struct {
		name      string
		userID    int64
		email     string
		modify    func(token string) string
		expectErr error
	}{
		{
			name:      "valid_access_token",
			userID:    testUserID,
			email:     testEmail,
			modify:    func(token string) string { return token },
			expectErr: nil,
		},
		{
			name:      "expired_access_token",
			userID:    testUserID,
			email:     testEmail,
			modify:    func(_ string) string { return string(expiredToken) },
			expectErr: domain.ErrTokenExpired,
		},
		{
			name:      "invalid_access_token",
			userID:    testUserID,
			email:     testEmail,
			modify:    func(token string) string { return token + "invalid" },
			expectErr: domain.ErrTokenInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			token, err := p.CreateToken(tt.userID, tt.email)
			require.NoError(t, err)

			token.Access = tt.modify(token.Access)
			payload, err := p.VerifyToken(token.Access)

			if tt.expectErr == nil {
				require.NoError(t, err)
				require.NotNil(t, payload)
				require.Equal(t, tt.userID, payload.UserID)
				require.Equal(t, tt.email, payload.Email)
			} else {
				require.ErrorIs(t, err, tt.expectErr)
				require.Nil(t, payload)
			}
		})
	}
}
