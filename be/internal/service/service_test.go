package service

import (
	"context"
	"errors"
	"testing"

	"github.com/escalopa/vego/internal/domain"
	"github.com/escalopa/vego/internal/service/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestService_GetOAuthRedirectURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		provider string
		wantURL  string
		wantErr  error
	}{
		{"valid_provider", "google", "http://redirect.url", nil},
		{"invalid_provider", "unknown", "", errors.New("invalid provider")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			op := mock.NewMockoauthProvider(ctrl)
			svc := New(nil, nil, op, nil, nil)

			op.EXPECT().GetRedirectURL(tt.provider).Return(tt.wantURL, tt.wantErr)
			url, err := svc.GetOAuthRedirectURL(tt.provider)
			require.Equal(t, tt.wantURL, url)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_RegisterUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		provider string
		code     string
		wantErr  error
	}{
		{"valid_registration", "google", "valid_code", nil},
		{"invalid_registration", "google", "invalid_code", errors.New("invalid code")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			op := mock.NewMockoauthProvider(ctrl)
			db := mock.NewMockdatabase(ctrl)
			up := mock.NewMockuserTokenProvider(ctrl)
			svc := New(db, nil, op, up, nil)

			user := &domain.User{Email: "test@example.com"}
			op.EXPECT().HandleCallback(gomock.Any(), tt.provider, tt.code).Return(user, tt.wantErr)
			if tt.wantErr == nil {
				db.EXPECT().CreateUser(gomock.Any(), user, tt.provider).Return(int64(1), nil)
				up.EXPECT().CreateToken(int64(1), user.Email).Return(&domain.Token{}, nil)
			}
			_, err := svc.RegisterUser(context.Background(), tt.provider, tt.code)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_AuthenticateUser(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		token   *domain.Token
		wantErr error
	}{
		{"valid_token", &domain.Token{Access: "valid_token"}, nil},
		{"invalid_token", &domain.Token{Access: "invalid_token"}, errors.New("invalid token")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mock.NewMockdatabase(ctrl)
			utp := mock.NewMockuserTokenProvider(ctrl)
			svc := New(db, nil, nil, utp, nil)

			payload := &domain.UserTokenPayload{UserID: 1}
			user := &domain.User{}
			utp.EXPECT().VerifyToken(tt.token.Access).Return(payload, tt.wantErr)
			if tt.wantErr == nil {
				db.EXPECT().GetUser(gomock.Any(), payload.UserID).Return(user, nil)
			}
			_, _, err := svc.AuthenticateUser(context.Background(), tt.token)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_CreateRoomToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		userID  int64
		roomID  string
		wantErr error
	}{
		{"valid_token", 1, "room1", nil},
		{"invalid_token", 1, "room1", errors.New("token creation failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			rtp := mock.NewMockroomTokenProvider(ctrl)
			svc := New(nil, nil, nil, nil, rtp)

			rtp.EXPECT().CreateToken(tt.userID, tt.roomID).Return("token", tt.wantErr)
			_, err := svc.CreateRoomToken(tt.userID, tt.roomID)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_AuthenticateWS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		token   string
		roomID  string
		wantErr error
	}{
		{"valid_token", "valid_token", "room1", nil},
		{"invalid_token", "invalid_token", "room1", errors.New("invalid token")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			db := mock.NewMockdatabase(ctrl)
			rtp := mock.NewMockroomTokenProvider(ctrl)
			svc := New(db, nil, nil, nil, rtp)

			payload := &domain.RoomTokenPayload{UserID: 1, RoomID: tt.roomID}
			user := &domain.User{}
			rtp.EXPECT().VerifyToken(tt.token).Return(payload, tt.wantErr)
			if tt.wantErr == nil {
				db.EXPECT().GetUser(gomock.Any(), payload.UserID).Return(user, nil)
			}
			_, err := svc.AuthenticateWS(context.Background(), tt.token, tt.roomID)
			require.Equal(t, tt.wantErr, err)
		})
	}
}

func TestService_HandleWS(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	h := mock.NewMockhub(ctrl)
	svc := New(nil, h, nil, nil, nil)

	user := &domain.User{}
	conn := &websocket.Conn{}
	roomID := "room1"

	h.EXPECT().Handle(user, roomID, conn).Times(1)
	svc.HandleWS(user, roomID, conn)
}
