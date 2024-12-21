package models

type Achievement struct {
	ID         uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Nama       string `json:"nama"`
	Pencapaian string `json:"pencapaian"`
	Link       string `json:"link"`
	Kategori   string `json:"kategori" validate:"required,oneof='International' 'National' 'Campus'"`
	Foto       string `json:"foto"`
	ImageURL   string `json:"image_url"`
}

func (Achievement) TableName() string {
	return "achievements"
}
