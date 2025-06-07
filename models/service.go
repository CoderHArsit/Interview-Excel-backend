package models

import "gorm.io/gorm"

func InitExpertRepo(db *gorm.DB) IExpert {
	return &expertRepo{
		DB: db,
	}
}
