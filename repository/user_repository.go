package repository

import (
	"context"

	"github.com/pratikgagare03/feedback/database"
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = database.Setup()
	if err != nil {
		panic(err)
	}
}

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.UserInput) error
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.UserInput) error
	DeleteUser(ctx context.Context, userID string) error
	GetUsers(tagcontains string) ([]models.User, error)
}

type postgresUserRepository struct {
	postgresDb *gorm.DB
}

// DeleteUser implements UserRepository.
func (p *postgresUserRepository) DeleteUser(ctx context.Context, userID string) error {
	panic("unimplemented")
}

// FindUserByID implements UserRepository.
func (p *postgresUserRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	panic("unimplemented")
}

// GetUsers implements UserRepository.
func (p *postgresUserRepository) GetUsers(tagcontains string) ([]models.User, error) {
	panic("unimplemented")
}

// InsertUser implements UserRepository.
func (p *postgresUserRepository) InsertUser(ctx context.Context, user *models.UserInput) error {
	var newUser = models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}
	res := Db.Create(&newUser)
	return res.Error
}

// UpdateUser implements UserRepository.
func (p *postgresUserRepository) UpdateUser(ctx context.Context, user *models.UserInput) error {
	panic("unimplemented")
}

func newPostgresUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{
		postgresDb: db,
	}
}
func GetUserRepository() UserRepository {
	return newPostgresUserRepository(Db)
}
