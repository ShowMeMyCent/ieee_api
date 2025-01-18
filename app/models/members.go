package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Members struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"`
	Role     Role   `json:"role" gorm:"type:enum('superAdmin' ,'admin', 'member', 'it'); default:'member'"`
}
type Role string

const (
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
	RoleIT     Role = "it"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleMember, RoleIT:
		return true
	}
	return true
}

func (m *Members) BeforeSave(*gorm.DB) (err error) {
	if m.Name == "" || m.Email == "" || m.Password == "" {
		return errors.New("name, email, and password are required")
	}

	if !m.Role.IsValid() {
		m.Role = RoleMember
	}

	if m.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		m.Password = string(hashedPassword)
	}
	return nil
}
