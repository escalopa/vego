package db

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/escalopa/vego/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *DB {
	t.Helper()

	testDBFile := fmt.Sprintf("/tmp/test_db_%s.db", uuid.NewString())
	db, err := New(testDBFile)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.conn.Close())
		require.NoError(t, os.Remove(testDBFile))
	})

	return db
}

func TestDBUserMethods(t *testing.T) {
	db := setupTestDB(t)

	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func() (int64, error)
		test      func(userID int64) error
		expectErr error
	}{
		{
			name: "create_and_get_user",
			setup: func() (int64, error) {
				user := &domain.User{
					Name:   "John Doe",
					Email:  "john.doe@example.com",
					Avatar: "avatar_url",
				}
				return db.CreateUser(ctx, user, "google")
			},
			test: func(userID int64) error {
				user, err := db.GetUser(ctx, userID)
				if err != nil {
					return err
				}
				require.Equal(t, "John Doe", user.Name)
				require.Equal(t, "john.doe@example.com", user.Email)
				require.Equal(t, "avatar_url", user.Avatar)
				return nil
			},
			expectErr: nil,
		},
		{
			name: "get_user_not_found",
			setup: func() (int64, error) {
				return 9999, nil // Non-existent user ID
			},
			test: func(userID int64) error {
				_, err := db.GetUser(ctx, userID)
				return err
			},
			expectErr: domain.ErrDBUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := tt.setup()
			require.NoError(t, err)

			err = tt.test(userID)
			if tt.expectErr != nil {
				require.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
