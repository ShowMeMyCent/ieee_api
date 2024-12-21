package models

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

func (News) TableName() string {
	return "news"
}
