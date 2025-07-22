package controllers

import (
	"fmt"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)


func CreateRazorpayOrderHandler(c *gin.Context) {
	var req struct {
		SlotID uint `json:"slot_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	availabilityRepo := models.InitAvailabilitySlotRepo(config.DB)
	expertRepo := models.InitExpertRepo(config.DB)

	slot, err := availabilityRepo.GetByID(req.SlotID)
	if err != nil || slot.IsBooked {
		c.JSON(http.StatusNotFound, gin.H{"error": "slot not available"})
		return
	}

	expert, err := expertRepo.GetByID(uint64(slot.ExpertID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "expert not found"})
		return
	}

	// Calculate amount
	platformFee := int(0.10 * float64(expert.FeesPerSession))
	totalAmount := expert.FeesPerSession + platformFee

	// Razorpay needs amount in paise
	amountInPaise := totalAmount * 100

	// Create order params
	data := map[string]interface{}{
		"amount":          amountInPaise,
		"currency":        "INR",
		"receipt":         fmt.Sprintf("receipt_slot_%d", slot.ID),
		"payment_capture": 1,
		"notes": map[string]interface{}{
			"slot_id":    slot.ID,
			"expert_id":  expert.ID,
			"student_id": c.GetUint("studentID"),
		},
	}

	order, err := config.RazorpayClient.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":     order["id"],
		"amount":       totalAmount,
		"currency":     "INR",
		"razorpay_key": os.Getenv("RAZORPAY_KEY"),
		"slot_id":      slot.ID,
	})
}
