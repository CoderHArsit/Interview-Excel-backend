package models

import "time"

type AvailabilitySlot struct {
	ID        uint      `gorm:"primaryKey"`
	ExpertID  uint      `gorm:"not null;index"` // Foreign Key to Expert
	Date      time.Time `gorm:"not null;index"` // Date of the slot (e.g., 2025-06-08)
	StartTime time.Time `gorm:"not null"`       // Start time (e.g., 14:00)
	EndTime   time.Time `gorm:"not null"`       // End time (e.g., 15:00)
	IsBooked  bool      `gorm:"default:false"`  // true if already booked
	CreatedAt time.Time
	UpdatedAt time.Time
}
