package entity

type User struct {
	ID
	UUID     string `gorm:"type:varchar(100); not null" json:"uuid"`
	NickName string `gorm:"type:varchar(100); not null" json:"nick_name"`
	BaseAt
}
