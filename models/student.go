package models

import (
	"gorm.io/gorm"
	"time"
)

type Student struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	FullName  string    `gorm:"not null" json:"full_name"`
	Picture   string    `json:"picture"`
	Email     string    `gorm:"not null;uniqueIndex" json:"email"`
	Phone     string    `gorm:"not null;uniqueIndex" json:"phone"`
	Password  string    `gorm:"not null" json:"-"` // stored as hashed
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
