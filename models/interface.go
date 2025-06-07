package models

import "gorm.io/gorm"

type IExpert interface {
	Create(s *Expert) error
	CreateWithTx(tx *gorm.DB, s *Expert) error
	GetAllExpert(provider *Expert) (*[]Expert, error)
	GetByID(providerID uint64) (*Expert, error)
	GetWithTx(tx *gorm.DB, where *Expert) (*Expert, error)
	Update(where *Expert, a *Expert) error
	UpdateWithTx(tx *gorm.DB, where *Expert, a *Expert) error
	Delete(where uint64) error
}
