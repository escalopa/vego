package db

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/escalopa/vego/internal/domain"
	_ "github.com/mattn/go-sqlite3" // sqlite3 driver
)

type DB struct {
	conn *sql.DB
}

func New(filepath string) (*DB, error) {
	conn, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}

	const query = `
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			avatar TEXT NOT NULL,
			provider TEXT NOT NULL
		);
		
		CREATE UNIQUE INDEX IF NOT EXISTS idx_email_provider ON users (email, provider); -- ensure email+provider is unique
	`
	_, err = conn.Exec(query)
	if err != nil {
		return nil, err
	}

	return &DB{conn: conn}, nil
}

func (db *DB) GetUser(_ context.Context, userID int64) (*domain.User, error) {
	const query = `
		SELECT user_id,
		       name,
			   email,
			   avatar
		FROM users
		WHERE user_id = $1
	`

	row := db.conn.QueryRow(query, userID)

	var res domain.User
	err := row.Scan(&res.UserID, &res.Name, &res.Email, &res.Avatar)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrDBUserNotFound
		}
		log.Printf("db.GetUser: %v", err)
		return nil, domain.ErrDBQuery
	}

	return &res, nil
}

func (db *DB) CreateUser(_ context.Context, user *domain.User, provider string) (int64, error) {
	const query = `
		INSERT INTO users (name, email, avatar, provider)
		VALUES ($1, $2, $3, $4) ON CONFLICT DO
		UPDATE SET name = $1, avatar = $3
		RETURNING user_id
	`

	row := db.conn.QueryRow(query, user.Name, user.Email, user.Avatar, provider)

	var userID int64
	err := row.Scan(&userID)
	if err != nil {
		log.Printf("db.CreateUser: %v", err)
		return 0, domain.ErrDBQuery
	}

	return userID, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}
