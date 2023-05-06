package entity

import (
	"time"
)

type User struct {
	Acct      string    `gorm:"type:varchar(255);primaryKey;not null"`
	Pwd       string    `gorm:"type:varchar(255);not null"`
	FullName  string    `gorm:"type:varchar(255);not null;column:fullname"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt time.Time `gorm:"type:timestamp with time zone;not null;default:now()"`
}
