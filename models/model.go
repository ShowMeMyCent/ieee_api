package models

import (
	"backend/utils/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Paper struct {
	ID            uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Judul         string `json:"judul"`
	Abstrak       string `json:"abstrak"`
	Link          string `json:"link"`
	FilePaper     string `json:"file_paper"`
	Author        string `json:"author"`
	TanggalTerbit string `json:"tanggal_terbit"`
}

type Activities struct {
	ID       uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Title    string `json:"title"`
	Tanggal  string `json:"tanggal"`
	Gambar   string `json:"gambar"`
	ImageURL string `json:"image_url"`
}

type News struct {
	ID          uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Title       string `json:"title" validate:"required"`
	Kategori    string `json:"kategori" validate:"required,oneof=News Event"`
	Thumbnail   string `json:"thumbnail"`
	IsiKonten   string `json:"isi_konten" gorm:"type:json"`
	NamaPenulis string `json:"nama_penulis" validate:"required"`
	Link        string `json:"link"`
	ImageURL    string `json:"image_url"`
	Date        string `json:"date"`
}

type Achievement struct {
	ID         uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Nama       string `json:"nama"`
	Pencapaian string `json:"pencapaian"`
	Link       string `json:"link"`
	Kategori   string `json:"kategori" validate:"required,oneof='International' 'National' 'Campus'"`
	Foto       string `json:"foto"`
	ImageURL   string `json:"image_url"`
}

type User struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Role string `json:"role"`
	Code string `json:"code"`
}

func (u *User) LoginCheckAdmin(db *gorm.DB) (string, error) {
	var err error

	adminFromDB := User{}
	err = db.Model(&User{}).Where("code = ?", u.Code).Take(&adminFromDB).Error
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(adminFromDB.Code), []byte(u.Code))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(adminFromDB.ID, adminFromDB.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *User) SaveAdmin(db *gorm.DB) (*User, error) {
	// Turn kode into hash
	hashedKode, err := bcrypt.GenerateFromPassword([]byte(u.Code), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, err
	}
	u.Code = string(hashedKode)

	err = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}
