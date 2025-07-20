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
	GetAll() ([]Expert, error)
}

type IAvailabilitySlotRepo interface {
	GetAllByExpert(expertID uint) ([]AvailabilitySlot, error)
	GetAvailableByExpert(expertID uint) ([]AvailabilitySlot, error)
	GetByID(id uint) (*AvailabilitySlot, error)
	MarkAsBooked(id uint) error
	Delete(id uint) error
	Update(slot *AvailabilitySlot) error
	GetBookedByStudent(studentID uint) ([]AvailabilitySlot, error)
	GetBookedSlotsByExpert(expertID uint) ([]AvailabilitySlot, error)
}
