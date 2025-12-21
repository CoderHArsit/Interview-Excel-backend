// controllers/student_controller.go
package controllers

import (
	"encoding/json"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	logger "interviewexcel-backend-go/pkg/errors"
	"net/http"

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
		UserID:       uuid,
		Role:         user.Role,
		FullName:     user.FullName,
		Email:        user.Email,
		Phone:        safeString(user.Phone),
		Bio:          student.Bio,
		PreparingFor: student.PreparingFor,
		Sessions:     student.Sessions,
		Points:       student.Points,
		DateOfBirth:  student.DateOfBirth,
		City:         student.City,
		AboutMe:      student.AboutMe,
		Skills:       skills, 
	}

	c.JSON(http.StatusOK, resp)
}
func UpdateStudentProfile(c *gin.Context) {
	var request StudentProfile

	userID, exists := c.Get("user_uuid")
	if !exists {
		logger.Error("user not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userUUID, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		logger.Error("error binding request: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		logger.Error("failed to start transaction: ", tx.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	studentRepo := models.InitStudentRepo(tx)
	userRepo := models.InitUserRepo(tx)

	// -------- STUDENT UPDATE --------

	skillsJSON, err := json.Marshal(request.Skills)
	if err != nil {
		tx.Rollback()
		logger.Error("failed to marshal skills: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skills format"})
		return
	}

	studentUpdates := map[string]interface{}{
		"bio":           request.Bio,
		"preparing_for": request.PreparingFor,
		"date_of_birth": request.DateOfBirth,
		"city":          request.City,
		"about_me":      request.AboutMe,
		"skills":        datatypes.JSON(skillsJSON),
	}

	if err := studentRepo.UpdateByUserUUID(userUUID, studentUpdates); err != nil {
		tx.Rollback()
		logger.Error("error updating student: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student profile"})
		return
	}
	if err := userRepo.UpdateByUserUUID(userUUID, &models.User{
		FullName: request.FullName,
		Phone:    &request.Phone,
	}); err != nil {
		tx.Rollback()
		logger.Error("error updating user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
		return
	}

	if err := tx.Commit().Error; err != nil {
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
	experts, err := expertRepo.GetAllExpertsWithUserDetails()
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

	slots, err := availabilityRepo.GetAvailableByExpert((expertIDStr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch slots"})
		return
	}

	c.JSON(http.StatusOK, slots)
}

