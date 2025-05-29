package repository

import (
	"errors"

	"github.com/ashil-poojary/gostruct/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB  *gorm.DB
	RDB *gorm.DB
}

func NewUserRepository(db1 *gorm.DB, db2 *gorm.DB) *UserRepository {
	return &UserRepository{
		DB:  db1,
		RDB: db2,
	}
}

// CreateUser uses GORM's Create method
func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

// GetByEmail uses GORM Where + First
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := r.DB.Where("email = ?", email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetByID uses GORM First by primary key
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	user := &models.User{}
	err := r.DB.First(user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

// GetAll uses GORM Find
func (r *UserRepository) GetAll() ([]*models.User, error) {
	var users []*models.User
	err := r.DB.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
