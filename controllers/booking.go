package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	logger "interviewexcel-backend-go/pkg/errors"
)

// Booking flow (what the function must guarantee)
// Atomic requirements

// Slot must exist

// Slot must be available

// Slot must be booked once

// Google Meet link must be created

// If any step fails → rollback

// High-level flow
// BEGIN TX
//   ├─ Lock availability slot
//   ├─ Validate slot availability
//   ├─ Create Google Meet link
//   ├─ Create booking record
//   ├─ Mark slot as booked
// COMMIT

type BookSlotRequest struct {
	SlotID uint `json:"slot_id" binding:"required"`
	// In future: StudentID uint `json:"student_id"`
}

func BookSlotHandler(c *gin.Context) {
	slotID, _ := strconv.Atoi(c.Param("slot_id"))

	if err := BookExpertSlot(c, uint(slotID)); err != nil {
		logger.Error("error in booking slot: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Slot booked successfully"})
}
