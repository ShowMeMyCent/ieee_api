package models

type News struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `json:"title" validate:"required"`
	Kategori    string `json:"kategori" validate:"required,oneof=News Event"`
	Date        string `json:"date" validate:"required"`
	IsiKonten   string `json:"isi_konten" validate:"required"`
	NamaPenulis string `json:"nama_penulis" validate:"required"`
	Link        string `json:"link"`
	Thumbnail   string `json:"thumbnail"`
	ImageURL    string `json:"image_url"`
}
