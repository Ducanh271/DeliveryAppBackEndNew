package repository

import (
	"database/sql"
	"errors"
	"example.com/delivery-app/models"
	"fmt"
	"time"
)

type RefreshTokenRepository interface {
	Save(userID int64, token string, expiresAt time.Time) error
	GetByToken(token string) (*models.RefreshToken, error)
	DeleteByToken(token string) (bool, error)
	DeleteByUserID(userID int64) error
	UpdateToken(oldToken string, newToken string) error

	WithTx(tx *sql.Tx) RefreshTokenRepository
}

type refreshTokenRepo struct {
	db DBTX
}

func NewRefreshTokenRepository(db DBTX) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) WithTx(tx *sql.Tx) RefreshTokenRepository {
	return &refreshTokenRepo{db: tx}
}

func (r *refreshTokenRepo) Save(userID int64, token string, expiresAt time.Time) error {
	query := "INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES (?, ?, ?)"
	_, err := r.db.Exec(query, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepo) GetByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	query := "SELECT id, user_id, token, expires_at, created_at FROM refresh_tokens WHERE token = ?"
	err := r.db.QueryRow(query, token).Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get refresh token by token: %w", err)
	}
	return &rt, nil
}

func (r *refreshTokenRepo) DeleteByToken(token string) (bool, error) {
	query := "DELETE FROM refresh_tokens WHERE token = ?"
	_, err := r.db.Exec(query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, fmt.Errorf("Invalid Token: %w", err)
		}
		return false, fmt.Errorf("failed to delete refresh token by token: %w", err)

	}
	return true, nil
}

func (r *refreshTokenRepo) DeleteByUserID(userID int64) error {
	query := "DELETE FROM refresh_tokens WHERE user_id = ?"
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete refresh tokens by user ID: %w", err)
	}
	return nil
}

func (r *refreshTokenRepo) UpdateToken(oldToken string, newToken string) error {
	query := "UPDATE refresh_tokens SET token = ? WHERE token = ?"
	_, err := r.db.Exec(query, newToken, oldToken)
	if err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}
	return nil
}

func StartTokenCleanUp(db *sql.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for range ticker.C {
			_, err := db.Exec("DELETE FROM refresh_tokens WHERE expires_at < ?", time.Now())
			if err != nil {
				fmt.Printf("failed to clean up expired refresh tokens: %v\n", err)
			} else {
				fmt.Println("expired refresh tokens cleaned up successfully")
			}
		}
	}()

}
