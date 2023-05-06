package repository

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/wys1203/go-gorilla-example/users/entity"
)

type UserRepository interface {
	GetAll() ([]entity.User, error)
	SearchByFullname(fullname string) ([]entity.User, error)
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
