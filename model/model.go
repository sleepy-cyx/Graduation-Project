package model

type User struct {
	Id       uint   `gorm:"primary_key;auto_increment" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Salt     string `gorm:"not null" json:"salt"`
}
type Schedule struct {
	Id        uint   `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint   `gorm:"not null" json:"user_id"`
	Title     string `gorm:"not null" json:"title"`
	Comment   string `gorm:"not null" json:"comment"`
	StartTime string `gorm:"not null" json:"start_time"`
	EndTime   string `gorm:"not null" json:"end_time"`
}
