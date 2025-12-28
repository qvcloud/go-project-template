package entity

type User struct {
	ID
	UUID     string `gorm:"type:varchar(100); not null"`
	NickName string `gorm:"type:varchar(100); not null"`
	BaseAt
}
