package models

type Activities struct {
	ID       uint   `gorm:"primaryKey"`
	Title    string `json:"title"`
	Tanggal  string `json:"tanggal"`
	Gambar   string `json:"gambar"`
	ImageURL string `json:"image_url"`
}
