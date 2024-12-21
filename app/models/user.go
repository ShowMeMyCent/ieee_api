package models

type User struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Role string `json:"role"`
	Code string `json:"code"`
}

func (User) TableName() string {
	return "user"
}
