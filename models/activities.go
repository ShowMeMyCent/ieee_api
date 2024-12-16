package models

type Activities struct {
	ID       uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Title    string `json:"title"`
	Date     string `json:"tanggal"`
	Image    string `json:"gambar"`
	ImageURL string `json:"image_url"`
}
