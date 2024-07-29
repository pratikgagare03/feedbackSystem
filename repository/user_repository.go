package repository

import (
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
	InsertUser(user *models.User) error
	FindUserByID(userID string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(userID string) error
	GetUsers(tagcontains string) ([]models.User, error)
}

type postgresUserRepository struct {
	postgresDb *gorm.DB
}

// DeleteUser implements UserRepository.
func (p *postgresUserRepository) DeleteUser(userID string) error {
	panic("unimplemented")
}

// FindUserByID implements UserRepository.
func (p *postgresUserRepository) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	res := Db.First(&user, userID)
	return &user, res.Error
}

// GetUsers implements UserRepository.
func (p *postgresUserRepository) GetUsers(tagcontains string) ([]models.User, error) {
	panic("unimplemented")
}

// InsertUser implements UserRepository.
func (p *postgresUserRepository) InsertUser(user *models.User) error {
	res := Db.Create(&user)
	return res.Error
}

// UpdateUser implements UserRepository.
func (p *postgresUserRepository) UpdateUser(user *models.User) error {
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
