package models

import (
	"gorm.io/gorm"
	"time"
)

type Student struct {
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"uniqueIndex" json:"user_uuid"` // references User.UserUUID
	Bio          string    `json:"bio,omitempty"`
	Sessions     string    `json:"sessions"`
	Points       string    `json:"points"`
	PreparingFor string    `json:"preparing_for"`
	DateOfBirth  time.Time `json:"dob"`
	City         string    `json:"city"`
	AboutMe      string    `json:"about_me"`
	Skills       []string  `gorm:"type:json" json:"skills"` // JSON column for skills
}

type StudentRepo struct {
	db *gorm.DB
}

func InitStudentRepo(db *gorm.DB) *StudentRepo {
	return &StudentRepo{db: db}
}

// Create a new student
func (r *StudentRepo) Create(student *Student) error {
	return r.db.Create(student).Error
}

// Get by user UUID
func (r *StudentRepo) GetByUserUUID(uuid string) (*Student, error) {
	var student Student
	err := r.db.Where("user_id = ?", uuid).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// Update student (by user UUID)
func (r *StudentRepo) UpdateByUserUUID(uuid string, updates map[string]interface{}) error {
	return r.db.Model(&Student{}).
		Where("user_id = ?", uuid).
		Updates(updates).Error
}

// Delete student (by user UUID)
func (r *StudentRepo) DeleteByUserUUID(uuid string) error {
	return r.db.Where("user_id = ?", uuid).Delete(&Student{}).Error
}

// List all students
func (r *StudentRepo) ListAll() ([]Student, error) {
	var students []Student
	err := r.db.Find(&students).Error
	return students, err
}
