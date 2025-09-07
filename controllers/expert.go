package controllers

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	logger "interviewexcel-backend-go/pkg/errors"
)

func GetExpertBookingsHandler(c *gin.Context) {
	var (
		availabilityRepo = models.InitAvailabilitySlotRepo(config.DB)
	)
	expertIDInterface, exists := c.Get("expert_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	expertID := expertIDInterface.(uint)

	slots, err := availabilityRepo.GetBookedSlotsByExpert(expertID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch bookings"})
		return
	}

	c.JSON(http.StatusOK, slots)
}

func GetExpertProfile(c *gin.Context) {
	var (
		expertRepo = models.InitExpertRepo(config.DB)
		userRepo   = models.InitUserRepo(config.DB)
	)

	user_id, exists := c.Get("user_uuid")
	if !exists {
		logger.Error("User doesn't exist")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	uuid, ok := user_id.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user UUID"})
		return
	}

	expertResp, err := expertRepo.GetWithTx(config.DB, &models.Expert{UserID: uuid})
	if err != nil {
		logger.Error("error in getting expert: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Expert not found"})
		return
	}

	userResp, err := userRepo.GetByUUID(uuid)
	if err != nil {
		logger.Error("error in getting user: ", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	profile := ExpertProfile{
		UserID:   userResp.ID,
		UserUUID: userResp.UserUUID,
		FullName: userResp.FullName,
		Email:    userResp.Email,
		Picture:  userResp.Picture,
		Phone:    userResp.Phone,
		Role:     userResp.Role,

		Bio:                expertResp.Bio,
		Expertise:          expertResp.Expertise,
		ExperienceYears:    expertResp.ExperienceYears,
		ProfilePictureUrl:  expertResp.ProfilePictureUrl,
		FeesPerSession:     expertResp.FeesPerSession,
		Rating:             expertResp.Rating,             // if present
		TotalSessions:      expertResp.TotalSessions,      // if present
		VerificationStatus: expertResp.VerificationStatus, // if present
	}

	c.JSON(http.StatusOK, profile)

}

func UpdateExpertProfile(c *gin.Context) {
	var (
		request ExpertProfile
	)
	user_id, exists := c.Get("user_uuid")
	if !exists {
		logger.Error("User doesn't exist")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	uuid, ok := user_id.(string)
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

	tx := config.DB.Begin()
	if tx.Error != nil {
		logger.Error("error in starting transaction: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	expertRepo := models.InitExpertRepo(tx)
	userRepo := models.InitUserRepo(tx)

	expertRequest := &models.Expert{
		Bio:               request.Bio,
		Expertise:         request.Expertise,
		ExperienceYears:   request.ExperienceYears,
		Education:         request.Education,
		Languages:         request.Languages,
		Specializations:   request.Specializations,
		ProfilePictureUrl: request.ProfilePictureUrl,
		FeesPerSession:    request.FeesPerSession,
	}

	err = expertRepo.UpdateWithTx(tx, &models.Expert{UserID: uuid}, expertRequest)
	if err != nil {
		tx.Rollback()
		logger.Error("error in updating expert: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update expert profile"})
		return
	}

	err = userRepo.UpdateByUserUUID(uuid, &models.User{
		FullName: request.FullName,
		Phone:    request.Phone,
	})
	if err != nil {
		tx.Rollback()
		logger.Error("error in updating user: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user profile"})
	}
	err = tx.Commit().Error
	if err != nil {
		logger.Error("failed to commit transaction: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Transaction failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

}
