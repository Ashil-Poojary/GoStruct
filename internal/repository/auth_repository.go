package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/ashil-poojary/gostruct/internal/models"
	"gorm.io/gorm"
)

type AuthRepo struct {
	DB *gorm.DB
}

func NewAuthRepo(db1 *gorm.DB) *AuthRepo {
	return &AuthRepo{
		DB: db1,
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// Save refresh token to DB using GORM Create
func (r *AuthRepo) SaveRefreshToken(token string, userID int64, userAgent, ip string, expiresAt time.Time) error {
	hashedToken := hashToken(token)
	now := time.Now()

	rt := &models.RefreshToken{
		UserID:           userID,
		RefreshTokenHash: hashedToken,
		UserAgent:        userAgent,
		IPAddress:        ip,
		IssuedAt:         now,
		ExpiresAt:        expiresAt,
		Revoked:          false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	return r.DB.Create(rt).Error
}

// Validate refresh token using GORM Where + First
func (r *AuthRepo) ValidateRefreshToken(token string) (*models.RefreshToken, error) {
	hashedToken := hashToken(token)

	rt := &models.RefreshToken{}
	err := r.DB.Where("refresh_token_hash = ?", hashedToken).First(rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	if rt.Revoked {
		return nil, errors.New("token revoked")
	}

	if time.Now().After(rt.ExpiresAt) {
		return nil, errors.New("token expired")
	}

	return rt, nil
}

// Revoke refresh token by updating revoked flag with GORM Model + Update
func (r *AuthRepo) RevokeRefreshToken(token string) error {
	hashedToken := hashToken(token)

	// Update revoked flag to true
	return r.DB.Model(&models.RefreshToken{}).
		Where("refresh_token_hash = ?", hashedToken).
		Update("revoked", true).Error
}
