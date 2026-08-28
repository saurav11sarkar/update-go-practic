package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

var ErrInvalidRole = errors.New("role must be user or admin")

type User struct {
	ID       string    `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `json:"name" gorm:"not null"`
	Email    string    `json:"email" gorm:"not null;uniqueIndex"`
	Password string    `json:"password" gorm:"not null"`
	Role     Role      `json:"role" gorm:"type:varchar(20);not null;default:user;check:role IN ('user','admin')"`
	IsActive bool      `json:"is_active" gorm:"default:true"`
	CreateAt time.Time `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"update_at" gorm:"autoUpdateTime"`
}

func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Role == "" {
		u.Role = RoleUser
	}
	if u.Role != RoleUser && u.Role != RoleAdmin {
		return ErrInvalidRole
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hasshPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hasshPassword)
	return nil
}
