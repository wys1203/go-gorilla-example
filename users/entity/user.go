package entity

import (
	"time"
)

type User struct {
	Acct      string     `json:"acct" gorm:"type:varchar(255);primaryKey;not null"`
	Pwd       string     `json:"pwd" gorm:"type:varchar(255);not null"`
	FullName  string     `json:"fullname" gorm:"type:varchar(255);not null;column:fullname"`
	CreatedAt *time.Time `json:"createdAt,omitempty" gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" gorm:"type:timestamp with time zone;not null;default:now()"`
}
