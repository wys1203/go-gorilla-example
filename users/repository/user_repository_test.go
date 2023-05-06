package repository_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/wys1203/go-gorilla-example/users/entity"
	"github.com/wys1203/go-gorilla-example/users/repository"
)

type RepositorySuite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (suite *RepositorySuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, suite.mock, err = sqlmock.New()
	suite.NoError(err)

	suite.DB, err = gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	suite.NoError(err)
}

func (suite *RepositorySuite) AfterTest(_, _ string) {
	suite.NoError(suite.mock.ExpectationsWereMet())
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(RepositorySuite))
}

func (suite *RepositorySuite) TestGetAll() {
	rows := sqlmock.NewRows([]string{"acct", "pwd", "fullname", "created_at", "updated_at"})
	timeLayout := "2006-01-02T15:04:05Z"
	createdAt, _ := time.Parse(timeLayout, "2023-05-06T15:04:05Z")
	updatedAt, _ := time.Parse(timeLayout, "2023-05-06T15:04:05Z")

	rows.AddRow("acctid-1", "password1", "User One", createdAt, updatedAt).
		AddRow("acctid-2", "password2", "User Two", createdAt, updatedAt)

	suite.mock.ExpectQuery("^SELECT (.+) FROM \"users\"$").WillReturnRows(rows)

	userRepo := repository.NewUserRepository(suite.DB)

	users, err := userRepo.GetAll()
	assert.NoError(suite.T(), err)

	expectedUsers := []entity.User{
		{Acct: "acctid-1", Pwd: "password1", FullName: "User One", CreatedAt: &createdAt, UpdatedAt: &updatedAt},
		{Acct: "acctid-2", Pwd: "password2", FullName: "User Two", CreatedAt: &createdAt, UpdatedAt: &updatedAt},
	}

	assert.Equal(suite.T(), expectedUsers, users)
}

func (suite *RepositorySuite) TestSearchByFullname() {
	rows := sqlmock.NewRows([]string{"acct", "pwd", "fullname", "created_at", "updated_at"})
	timeLayout := "2006-01-02T15:04:05Z"
	createdAt, _ := time.Parse(timeLayout, "2023-05-06T15:04:05Z")
	updatedAt, _ := time.Parse(timeLayout, "2023-05-06T15:04:05Z")

	rows.
		AddRow("acctid-1", "password1", "John Doe", createdAt, updatedAt).
		AddRow("acctid-2", "password2", "Jane Doe", createdAt, updatedAt)

	suite.mock.ExpectQuery("SELECT (.+) FROM \"users\" WHERE fullname LIKE (.+)").
		WithArgs("%Doe%").
		WillReturnRows(rows)

	userRepo := repository.NewUserRepository(suite.DB)

	users, err := userRepo.SearchByFullname("Doe")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), users, 2)

	expectedUsers := []entity.User{
		{
			Acct:     "acctid-1",
			Pwd:      "password1",
			FullName: "John Doe",
		},
		{
			Acct:     "acctid-2",
			Pwd:      "password2",
			FullName: "Jane Doe",
		},
	}

	for i, user := range users {
		assert.Equal(suite.T(), expectedUsers[i].Acct, user.Acct)
		assert.Equal(suite.T(), expectedUsers[i].Pwd, user.Pwd)
		assert.Equal(suite.T(), expectedUsers[i].FullName, user.FullName)
	}

}
