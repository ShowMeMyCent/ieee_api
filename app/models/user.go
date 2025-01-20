package models

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement:true"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email" gorm:"unique"`
	Password string `json:"-" binding:"required"`
	Role     Role   `json:"role" gorm:"type:enum('super_admin' ,'admin', 'member', 'tim_it');default:'member'"`
}

type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleMember     Role = "member"
	RoleIT         Role = "tim_it"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleMember, RoleIT:
		return true
	default:
		return false
	}
}
