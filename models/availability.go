package models

import "time"

type AvailabilitySlot struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ExpertID  uint      `gorm:"not null;index" json:"expert_id"` // Foreign Key to Expert
	Expert    Expert    `gorm:"foreignKey:ExpertID" json:"-"`    // Optional: to preload Expert info if needed
	Date      time.Time `gorm:"not null;index" json:"date"`      // Date of the slot (e.g., 2025-06-08)
	StartTime time.Time `gorm:"not null" json:"start_time"`      // Start time (e.g., 14:00)
	EndTime   time.Time `gorm:"not null" json:"end_time"`        // End time (e.g., 15:00)
	IsBooked  bool      `gorm:"default:false" json:"is_booked"`  // true if already booked
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}


