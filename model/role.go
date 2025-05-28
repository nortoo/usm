package model

import "time"

type (
	// A Role represents a set of permissions and menus that can be assigned to users.
	Role struct {
		ID uint

		Name    string `gorm:"unique; index; not null"`
		Comment string `gorm:"not null; default: ''"`

		ApplicationID uint          `gorm:"not null"`
		Application   *Application  `gorm:"foreignKey:ApplicationID; references:ID; onDelete:CASCADE"`
		Menus         []*Menu       `gorm:"many2many:role_menus;"`
		Permissions   []*Permission `gorm:"many2many:role_permissions;"`

		CreatedAt time.Time
		UpdatedAt time.Time
	}
)

func init() {
	registerModel(new(Role))
}
