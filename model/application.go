package model

import (
	"time"
)

type (
	// Application acts as a tenant in the system, representing a distinct source.
	Application struct {
		ID        uint
		Name      string `gorm:"unique; not null; index; size:64; comment: application name"`
		APPID     string `gorm:"column:appid; not null; unique; index; size:32; comment: application ID"`
		SecretKey string `gorm:"not null; size:64"`
		Comment   string `gorm:"not null; default: ''"`

		State     int8 `gorm:"not null; default:1;comment: application state, 1 for active, 0 for inactive"`
		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func init() {
	registerModel(new(Application))
}
