package repository

import (
	"github.com/pratikgagare03/feedback/database"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() error {
	var err error
	db, err = database.Connect()
	if err != nil {
		logger.Logs.Fatal().Msgf("Error connecting to the database: %v", err)
	}
	return err
}

type UserRepository interface {
	InsertUser(user *models.User) error
	FindUserByID(userID uint) (*models.User, error)
	FindUserByEmail(email string) (*models.User, error)
	GetAllUsersByOffsetLimit(offset, limit int) ([]models.User, error)
	GetUsersCount() (int64, error)
	DeleteUserByEmail(email string) error
}

type postgresUserRepository struct {
	postgresdb *gorm.DB
}

// DeleteUserByEmail implements UserRepository.
func (p *postgresUserRepository) DeleteUserByEmail(email string) error {
	res := db.Unscoped().Delete(&models.User{}, "email = ?", email)
	return res.Error
}

// GetTotalUsers implements UserRepository.
func (p *postgresUserRepository) GetUsersCount() (int64, error) {
	var count int64
	res := db.Model(&models.User{}).Count(&count)
	return count, res.Error
}

// FindUserByEmail implements UserRepository.
func (p *postgresUserRepository) FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	res := db.First(&user, "email=?", email)
	return &user, res.Error
}

// FindUserByID implements UserRepository.
func (p *postgresUserRepository) FindUserByID(userID uint) (*models.User, error) {
	var user models.User
	user.ID = userID
	res := db.First(&user)
	return &user, res.Error
}

// GetUsers implements UserRepository.
func (p *postgresUserRepository) GetAllUsersByOffsetLimit(offset, limit int) ([]models.User, error) {
	var users []models.User
	res := db.Offset(offset).Limit(limit).Find(&users)
	return users, res.Error
}

// InsertUser implements UserRepository.
func (p *postgresUserRepository) InsertUser(user *models.User) error {
	res := db.Create(&user)
	return res.Error
}

func newPostgresUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{
		postgresdb: db,
	}
}
func GetUserRepository() UserRepository {
	return newPostgresUserRepository(db)
}
