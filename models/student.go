package models

import (
	"gorm.io/gorm"
	"time"
)

type Student struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex" json:"user_id"` // 1-to-1 with User
	Bio       string    `json:"bio,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StudentRepo struct {
	db *gorm.DB
}

func InitStudentRepo(db *gorm.DB) *StudentRepo {
	return &StudentRepo{db: db}
}

func (r *StudentRepo) Create(student *Student) error {
	return r.db.Create(student).Error
}

func (r *StudentRepo) GetByEmail(email string) (*Student, error) {
	var student Student
	err := r.db.Where("email = ?", email).First(&student).Error
	return &student, err
}
