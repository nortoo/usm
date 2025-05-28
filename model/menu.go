package model

import "time"

type (
	// Menu represents a menu item in the system.
	Menu struct {
		ID uint

		// ParentID is the ID of the parent menu item;
		//`0` means it is a top-level menu
		ParentID int64 `gorm:"not null; default: 0"`

		Name    string `gorm:"unique; not null"`
		Path    string `gorm:"unique; not null"`
		Comment string `gorm:"not null; default: ''"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func init() {
	registerModel(new(Menu))
}
