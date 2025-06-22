package model

import "time"

type (
	Group struct {
		ID uint

		Name      string  `gorm:"unique; not null;"`
		Comment   string  `gorm:"not null; default: ''"`
		Users     []*User `gorm:"many2many:user_groups;"`
		IsDefault bool    `gorm:"default:false"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func init() {
	registerModel(new(Group))
}
