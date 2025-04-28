package model

type User struct {
	Id         uint   `gorm:"primary_key;auto_increment" json:"id"`
	Username   string `gorm:"unique" json:"username"`
	Password   string `gorm:"not null" json:"password"`
	Salt       string `gorm:"not null" json:"salt"`
	Email      string `gorm:"unique" json:"email"`
	PictureUrl string `gorm:"not null" json:"picture_url"`
	NeedNotice bool   `gorm:"not null" json:"need_notice"`
}
type Schedule struct {
	Id        uint   `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	Title     string `gorm:"not null" json:"title"`
	Location  string `gorm:"not null" json:"location"`
	Comment   string `gorm:"not null" json:"comment"`
	StartTime string `gorm:"not null" json:"start_time"`
	EndTime   string `gorm:"not null" json:"end_time"`
	Type      uint   `gorm:"not null" json:"type"`
}
