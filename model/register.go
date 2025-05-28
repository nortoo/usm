package model

import "gorm.io/gorm"

var models = make([]any, 0)

func registerModel(model any) {
	models = append(models, model)
}

// RegisterModels registers all models with the provided GORM database instance.
func RegisterModels(db *gorm.DB) error {
	return db.AutoMigrate(models...)
}
