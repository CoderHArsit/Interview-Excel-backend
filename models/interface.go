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

type IPaymentRepo interface {
	Create(payment *Payment) error
	GetByOrderID(orderID string) (*Payment, error)
}

type IStudent interface {
	Create(student *Student) error
	GetByEmail(email string) (*Student, error)
}

type IUser interface {
	InitUserRepo(db *gorm.DB) *UserRepo
	Create(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	ListAll() ([]User, error)
}
