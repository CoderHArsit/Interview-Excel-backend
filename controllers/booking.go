package controllers

import (
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BookSlotRequest struct {
	SlotID uint `json:"slot_id" binding:"required"`
	// In future: StudentID uint `json:"student_id"`
}

func BookSlot(c *gin.Context) {
	var request BookSlotRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slotRepo := models.InitAvailabilitySlotRepo(config.DB)

	// 1. Check if the slot exists
	slot, err := slotRepo.GetByID(request.SlotID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slot not found"})
		return
	}

	// 2. Check if already booked
	if slot.IsBooked {
		c.JSON(http.StatusConflict, gin.H{"error": "Slot already booked"})
		return
	}

	// 3. Mark as booked
	err = slotRepo.MarkAsBooked(request.SlotID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to book slot"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slot booked successfully"})
}
