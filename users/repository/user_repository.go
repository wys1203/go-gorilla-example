package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/wys1203/go-gorilla-example/users/entity"
)

type UserRepository interface {
	GetAll() ([]entity.User, error)
	SearchByFullname(fullname string) ([]entity.User, error)
	GetByAcct(acct string) (*entity.User, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepositoryImpl) SearchByFullname(fullname string) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Where("fullname LIKE ?", fmt.Sprintf("%%%s%%", fullname)).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepositoryImpl) GetByAcct(acct string) (*entity.User, error) {
	var user entity.User

	// Use the GORM library to search for a user with a matching account ID
	err := r.db.Where("acct = ?", acct).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}
