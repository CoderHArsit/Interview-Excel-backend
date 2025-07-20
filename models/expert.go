package models

import (
	"time"

	"gorm.io/gorm"
)

type Expert struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	FullName          string    `gorm:"not null" json:"full_name"`
	Email             string    `json:"email"`
	Picture           string    `json:"picture"`
	Password          string    `gorm:"not null" json:"-"`
	Phone             string    `json:"phone,omitempty" gorm:"unique;not null"`
	Expertise         string    `gorm:"not null" json:"expertise"`
	Bio               string    `json:"bio,omitempty"`
	ExperienceYears   int       `gorm:"not null" json:"experience_years"`
	ProfilePictureUrl string    `json:"profile_picture_url,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
type expertRepo struct {
	DB *gorm.DB
}

func (e *expertRepo) Create(s *Expert) error {

	return e.CreateWithTx(e.DB, s)
}
func (e *expertRepo) CreateWithTx(tx *gorm.DB, s *Expert) error {

	err := tx.Create(s).Error
	if err != nil {
		return err
	}
	return nil
}
func (e *expertRepo) GetAllExpert(ex *Expert) (*[]Expert, error) {
	experts := []Expert{}
	err := e.DB.Where(ex).Find(&experts).Error
	if err != nil {
		return nil, err
	}
	return &experts, nil
}

func (e *expertRepo) GetByID(ID uint64) (*Expert, error) {
	var expert Expert
	err := e.DB.First(&expert, ID).Error
	if err != nil {
		return nil, err
	}
	return &expert, nil
}

func (e *expertRepo) GetWithTx(tx *gorm.DB, where *Expert) (*Expert, error) {
	var expert Expert
	err := tx.Where(where).First(&expert).Error
	if err != nil {
		return nil, err
	}
	return &expert, nil
}

func (e *expertRepo) Update(where *Expert, update *Expert) error {
	return e.UpdateWithTx(e.DB, where, update)
}

func (e *expertRepo) UpdateWithTx(tx *gorm.DB, where *Expert, update *Expert) error {
	err := tx.Model(&Expert{}).Where(&where).Updates(&update).Error
	if err != nil {
		return err
	}
	return nil
}

func (e *expertRepo) Delete(id uint64) error {
	err := e.DB.Delete(&Expert{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *expertRepo) GetAll() ([]Expert, error) {
	var experts []Expert
	err := r.DB.Find(&experts).Error
	return experts, err
}


