// controllers/student_controller.go
package controllers

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetStudentProfile(c *gin.Context) {
	studentID, exists := c.Get("student_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var student models.Student
	if err := config.DB.First(&student, studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func UpdateStudentProfile(c *gin.Context) {
	studentID, exists := c.Get("student_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input struct {
		Name         string `json:"name"`
		Bio          string `json:"bio"`
		ProfileImage string `json:"profile_image"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := config.DB.Model(&models.Student{}).
		Where("id = ?", studentID).
		Updates(models.Student{
			FullName: input.Name,
			Bio:      input.Bio,
			Picture:  input.ProfileImage,
		}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetAllExpertsHandler(c *gin.Context) {
	var (
		expertRepo = models.InitExpertRepo(config.DB)
	)
	experts, err := expertRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch experts"})
		return
	}

	c.JSON(http.StatusOK, experts)
}

func GetAvailableSlotsForExpertHandler(c *gin.Context) {
	var (
		availabilityRepo = models.InitAvailabilitySlotRepo(config.DB)
	)
	expertIDStr := c.Param("id")
	expertID, err := strconv.Atoi(expertIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expert id"})
		return
	}

	slots, err := availabilityRepo.GetAvailableByExpert(uint(expertID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch slots"})
		return
	}

	c.JSON(http.StatusOK, slots)
}

func BookAvailabilitySlotHandler(c *gin.Context) {
	// Extract student ID from context
	var (
		availabilityRepo = models.InitAvailabilitySlotRepo(config.DB)
	)
	studentIDVal, exists := c.Get("studentID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	studentID := studentIDVal.(uint)

	// Parse body
	var req BookSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Get slot
	slot, err := availabilityRepo.GetByID(req.SlotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "slot not found"})
		return
	}

	if slot.IsBooked {
		c.JSON(http.StatusBadRequest, gin.H{"error": "slot already booked"})
		return
	}

	// Book slot
	slot.IsBooked = true
	slot.StudentID = &studentID
	if err := availabilityRepo.Update(slot); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to book slot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "slot booked successfully", "slot": slot})
}


func GetStudentBookingsHandler(c *gin.Context) {
	var(
		availabilityRepo = models.InitAvailabilitySlotRepo(config.DB)
	)
	studentIDInterface, exists := c.Get("student_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	studentID := studentIDInterface.(uint)

	slots, err := availabilityRepo.GetBookedByStudent(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, slots)
}
