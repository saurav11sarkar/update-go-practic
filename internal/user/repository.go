package user

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type Repository interface {
	GetUserByEmail(email string) (*User, error)
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where(&User{Email: email}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) CreateUser(user *User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) ||
			(errors.As(result.Error, &pgErr) && pgErr.Code == "23505") {
			return ErrUserAlreadyExists
		}
		return result.Error
	}
	return nil
}

func (r *repository) GetUserByID(id string) (*User, error) {
	var user User
	if err := r.db.Where(&User{ID: id}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
