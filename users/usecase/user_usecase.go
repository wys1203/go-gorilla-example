package usecase

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt"

	"github.com/wys1203/go-gorilla-example/users/entity"
	"github.com/wys1203/go-gorilla-example/users/repository"
)

type UserUsecase interface {
	GetAll(page int, size int, sortBy string, order string) ([]entity.User, error)
	SearchUsers(fullname string) ([]entity.User, error)
	GetUserByAcct(acct string) (*entity.User, error)
	CreateUser(user *entity.User) (*entity.User, error)
	Login(acct, pwd string) (string, error)
	Delete(acct string) error
	Update(acct string, user entity.User) error
}

type UserUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &UserUsecaseImpl{userRepo: userRepo}
}

func (u *UserUsecaseImpl) GetAll(page int, size int, sortBy string, order string) ([]entity.User, error) {
	return u.userRepo.GetAll(page, size, sortBy, order)
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

func (u *UserUsecaseImpl) Login(acct, pwd string) (string, error) {
	user, err := u.userRepo.GetByAcct(acct)
	if err != nil {
		return "", err
	}

	// Verify that the password matches the one in the database
	if user.Pwd != pwd {
		return "", fmt.Errorf("invalid password")
	}

	privKeyBytes, err := ioutil.ReadFile("private_key.pem")
	if err != nil {
		return "", err
	}

	privKeyPEM, _ := pem.Decode(privKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privKeyPEM.Bytes)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": user.Acct,
	})

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (u *UserUsecaseImpl) Delete(acct string) error {
	return u.userRepo.Delete(acct)
}

func (u *UserUsecaseImpl) Update(acct string, user entity.User) error {
	return u.userRepo.Update(acct, user)
}
