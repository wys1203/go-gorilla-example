package repository

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/wys1203/go-gorilla-example/users/entity"
)

type UserRepository interface {
	GetAll(page int, size int, sortBy string, order string) ([]entity.User, error)
	SearchByFullname(fullname string) ([]entity.User, error)
	GetByAcct(acct string) (*entity.User, error)
	Create(user *entity.User) (*entity.User, error)
	Delete(acct string) error
	Update(acct string, user entity.User) error
	UpdateFullname(acct string, fullname string) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) GetAll(page int, size int, sortBy string, order string) ([]entity.User, error) {
	if sortBy == "" {
		sortBy = "created_at"
	}

	if order == "" {
		order = "asc"
	}

	if page <= 0 {
		page = 1
	}

	switch {
	case size > 100:
		size = 100
	case size <= 0:
		size = 10
	}

	offset := (page - 1) * size

	var users []entity.User
	err := r.db.Offset(offset).Limit(size).Order(fmt.Sprintf("%s %s", sortBy, order)).Find(&users).Error
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

func (r *UserRepositoryImpl) Create(user *entity.User) (*entity.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepositoryImpl) Delete(acct string) error {
	return r.db.Delete(&entity.User{}, "acct = ?", acct).Error
}

func (r *UserRepositoryImpl) Update(acct string, user entity.User) error {
	// Update attributes with `struct`, will only update non-zero fields
	res := r.db.Model(&entity.User{}).Where("acct = ?", acct).Updates(user)

	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepositoryImpl) UpdateFullname(acct string, fullname string) error {
	res := r.db.Model(&entity.User{}).Where("acct = ?", acct).Update("fullname", fullname)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
