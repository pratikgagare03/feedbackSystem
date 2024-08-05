package repository

import (
	"github.com/pratikgagare03/feedback/database"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	var err error
	Db, err = database.Connect()
	if err != nil {
		logger.Logs.Fatal().Msgf("Error connecting to the database: %v", err)
	}
}

type UserRepository interface {
	InsertUser(user *models.User) error
	FindUserByID(userID uint) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	DeleteUser(userID string) error
	GetAllUsersByOffsetLimit(offset, limit int) ([]models.User, error)
	GetUsersCount() (int64, error)
}

type postgresUserRepository struct {
	postgresDb *gorm.DB
}

// GetTotalUsers implements UserRepository.
func (p *postgresUserRepository) GetUsersCount() (int64, error) {
	var count int64
	res := Db.Model(&models.User{}).Count(&count)
	return count, res.Error
}

// FindUserByEmail implements UserRepository.
func (p *postgresUserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	res := Db.First(&user, "email=?", email)
	return &user, res.Error
}

// DeleteUser implements UserRepository.
func (p *postgresUserRepository) DeleteUser(userID string) error {
	panic("unimplemented")
}

// FindUserByID implements UserRepository.
func (p *postgresUserRepository) FindUserByID(userID uint) (*models.User, error) {
	var user models.User
	user.ID = userID
	res := Db.First(&user)
	return &user, res.Error
}

// GetUsers implements UserRepository.
func (p *postgresUserRepository) GetAllUsersByOffsetLimit(offset, limit int) ([]models.User, error) {
	var users []models.User
	res := Db.Offset(offset).Limit(limit).Find(&users)
	return users, res.Error
}

// InsertUser implements UserRepository.
func (p *postgresUserRepository) InsertUser(user *models.User) error {
	res := Db.Create(&user)
	return res.Error
}


func newPostgresUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{
		postgresDb: db,
	}
}
func GetUserRepository() UserRepository {
	return newPostgresUserRepository(Db)
}
