package entities

type User struct {
	ID
	Uuid     string `gorm:"type:varchar(100); not null"`
	NickName string `gorm:"type:varchar(100); not null"`
	BaseAt
}
