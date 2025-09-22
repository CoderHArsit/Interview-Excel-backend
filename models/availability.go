package models

import (
	"time"

	"gorm.io/gorm"
)

type AvailabilitySlot struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	ExpertID string `gorm:"not null;index" json:"expert_id"` // references Expert.UserID
	Expert   Expert `gorm:"foreignKey:ExpertID;references:UserID" json:"-"`

	Date      time.Time `gorm:"not null;index" json:"date"`
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	IsBooked  bool      `gorm:"default:false" json:"is_booked"`
	StudentID *uint     `gorm:"index" json:"student_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type availabilitySlotRepo struct {
	DB *gorm.DB
}

func InitAvailabilitySlotRepo(db *gorm.DB) IAvailabilitySlotRepo {
	return &availabilitySlotRepo{
		DB: db,
	}
}
func (r *availabilitySlotRepo) CreateAvailabilitySlot(availability []AvailabilitySlot) error {
	return r.DB.Create(&availability).Error
}

// Get all slots for an expert
func (r *availabilitySlotRepo) GetAllByExpert(expertID string) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	err := r.DB.Where("expert_id = ?", expertID).Find(&slots).Error
	return slots, err
}

// Get all available (not booked) slots
func (r *availabilitySlotRepo) GetAvailableByExpert(expertID string) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	err := r.DB.Where("expert_id = ? AND is_booked = false AND date >= ?", expertID, time.Now()).
		Order("date ASC, start_time ASC").
		Find(&slots).Error
	return slots, err
}

// Get slot by ID
func (r *availabilitySlotRepo) GetByID(id uint) (*AvailabilitySlot, error) {
	var slot AvailabilitySlot
	err := r.DB.First(&slot, id).Error
	return &slot, err
}

// Mark slot as booked
func (r *availabilitySlotRepo) MarkAsBooked(id uint) error {
	return r.DB.Model(&AvailabilitySlot{}).Where("id = ?", id).Update("is_booked", true).Error
}

// Delete a slot
func (r *availabilitySlotRepo) Delete(id uint) error {
	return r.DB.Delete(&AvailabilitySlot{}, id).Error
}

// Update a slot (useful for admin panel or expert updates)
func (r *availabilitySlotRepo) Update(slot *AvailabilitySlot) error {
	return r.DB.Save(slot).Error
}

func (r *availabilitySlotRepo) GetBookedByStudent(studentID uint) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	err := r.DB.
		Where("student_id = ? AND is_booked = true AND date >= ?", studentID, time.Now()).
		Order("date ASC").
		Find(&slots).Error
	return slots, err
}

func (r *availabilitySlotRepo) GetBookedSlotsByExpert(expertID uint) ([]AvailabilitySlot, error) {
	var slots []AvailabilitySlot
	err := r.DB.
		Where("expert_id = ? AND is_booked = true AND date >= ?", expertID, time.Now()).
		Order("date ASC, start_time ASC").
		Find(&slots).Error
	return slots, err
}
