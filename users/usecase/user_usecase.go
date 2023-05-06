package usecase

import (
	"github.com/wys1203/go-gorilla-example/users/entity"
	"github.com/wys1203/go-gorilla-example/users/repository"
)

type UserUsecase interface {
	GetAllUsers() ([]entity.User, error)
	SearchUsers(fullname string) ([]entity.User, error)
	GetUserByAcct(acct string) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
}

type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{userRepo: userRepo}
}

func (u *UserUsecaseImpl) GetAllUsers() ([]entity.User, error) {
	users, err := u.userRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserUsecaseImpl) SearchUsers(fullname string) ([]entity.User, error) {
	return u.userRepo.SearchByFullname(fullname)
}

func (u *UserUsecaseImpl) GetUserByAcct(acct string) (*entity.User, error) {
	return u.userRepo.GetByAcct(acct)
}

func (u *UserUsecaseImpl) CreateUser(user *entity.User) (*entity.User, error) {
	return u.userRepo.Create(user)
}
