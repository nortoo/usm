package model

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Username string `gorm:"unique; index; not null"`
	Password string `gorm:"not null"`

	Email  string   `gorm:"index; not null; default: ''"`
	Mobile string   `gorm:"index; not null; default: ''"`
	Roles  []*Role  `gorm:"many2many:user_roles;"`
	Groups []*Group `gorm:"many2many:user_groups;"`

	State int8 `gorm:"not null; default: 0"`
}

func init() {
	registerModel(new(User))
}
