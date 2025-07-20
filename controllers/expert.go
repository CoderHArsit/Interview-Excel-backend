package controllers

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetExpertBookingsHandler(c *gin.Context) {
	var(
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
