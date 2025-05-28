package model

import "time"

type (
	// A Permission defines a combination of a source and their action.
	Permission struct {
		ID uint

		Action   string `gorm:"index:,unique,composite:action_resource; size:64; not null'"`
		Resource string `gorm:"index:,unique,composite:action_resource; size:256; not null"`
		Comment  string `gorm:"not null; default: ''"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func init() {
	registerModel(new(Permission))
}
