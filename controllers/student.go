// controllers/student_controller.go
package controllers

import (
	"encoding/json"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	logger "interviewexcel-backend-go/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func safeString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}
func GetStudentProfile(c *gin.Context) {
	studentRepo := models.InitStudentRepo(config.DB)
	userRepo := models.InitUserRepo(config.DB)

	// Extract user_uuid from context
	userID, exists := c.Get("user_uuid")
	if !exists {
		logger.Error("User not exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uuid, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID"})
		return
	}

	// Fetch student details
	student, err := studentRepo.GetByUserUUID(uuid)
	if err != nil {
		logger.Error("error in getting student: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Fetch user details
	user, err := userRepo.GetByUUID(uuid)
	if err != nil {
		logger.Error("error in getting user: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 🔹 Unmarshal skills JSON into []string
	var skills []string
	if len(student.Skills) > 0 {
		if err := json.Unmarshal(student.Skills, &skills); err != nil {
			logger.Error("error unmarshalling skills: ", err)
			skills = []string{} // fallback empty
		}
	}

	// Merge response
	resp := StudentProfile{
		UserID:      uuid,
		Role:        user.Role,
		FullName:    user.FullName,
		Email:       user.Email,
		Phone:       safeString(user.Phone),
		Bio:         student.Bio,
		Sessions:    student.Sessions,
		Points:      student.Points,
		DateOfBirth: student.DateOfBirth,
		City:        student.City,
		AboutMe:     student.AboutMe,
		Skills:      skills, // ✅ now proper []string
	}

	c.JSON(http.StatusOK, resp)
}

func UpdateStudentProfile(c *gin.Context) {
	var (
		request StudentProfile
	)

	userID, exists := c.Get("user_uuid")
	if !exists {
		logger.Error("User doesn't exist")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uuid, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID"})
		return
	}

	err := c.ShouldBindJSON(&request)
	if err != nil {
		logger.Error("error in binding Request: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// start a transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		logger.Error("error in starting transaction: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// use tx for repos
	studentRepo := models.InitStudentRepo(tx)
	userRepo := models.InitUserRepo(tx)

	skillsJSON, _ := json.Marshal(request.Skills)

	studentRequest := &models.Student{
		Bio:          request.Bio,
		PreparingFor: request.PreparingFor,
		DateOfBirth:  request.DateOfBirth,
		City:         request.City,
		AboutMe:      request.AboutMe,
		Skills:       datatypes.JSON(skillsJSON), // ✅ direct cast
	}

	// first update student
	err = studentRepo.UpdateByUserUUID(uuid, studentRequest)
	if err != nil {
		tx.Rollback()
		logger.Error("error in updating student: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student profile"})
		return
	}

	// then update user
	err = userRepo.UpdateByUserUUID(uuid, &models.User{
		FullName: request.FullName,
		Phone:    &request.Phone,
	})
	if err != nil {
		tx.Rollback()
		logger.Error("error in updating user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	// commit if both succeeded
	err = tx.Commit().Error
	if err != nil {
		logger.Error("failed to commit transaction: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
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

func PreviewSlotForPaymentHandler(c *gin.Context) {
	var req struct {
		SlotID uint `json:"slot_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	availabilityRepo := models.InitAvailabilitySlotRepo(config.DB)

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

	// Get expert
	expertRepo := models.InitExpertRepo(config.DB)
	expert, err := expertRepo.GetByID(uint64(slot.ExpertID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "expert not found"})
		return
	}

	// Platform fee logic
	platformFee := int(0.2 * float64(expert.FeesPerSession)) // 10%
	totalAmount := expert.FeesPerSession + platformFee

	c.JSON(http.StatusOK, gin.H{
		"expert": gin.H{
			"domain":           expert.Expertise,
			"fees_per_session": expert.FeesPerSession,
		},
		"slot": gin.H{
			"slot_id":    slot.ID,
			"day":        slot.Date,
			"start_time": slot.StartTime,
			"end_time":   slot.EndTime,
		},
		"platform_fee": platformFee,
		"total_amount": totalAmount,
	})
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
	var (
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

