package controllers

import (
	"context"
	"errors"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"gorm.io/gorm/clause"
)

func CreateGoogleMeetLink(
	ctx context.Context,
	client *calendar.Service,
	start, end time.Time,
) (string, error) {

	event := &calendar.Event{
		Summary: "InterviewExcel Expert Session",
		Start: &calendar.EventDateTime{
			DateTime: start.Format(time.RFC3339),
			TimeZone: "Asia/Kolkata",
		},
		End: &calendar.EventDateTime{
			DateTime: end.Format(time.RFC3339),
			TimeZone: "Asia/Kolkata",
		},
		ConferenceData: &calendar.ConferenceData{
			CreateRequest: &calendar.CreateConferenceRequest{
				RequestId: uuid.New().String(),
				ConferenceSolutionKey: &calendar.ConferenceSolutionKey{
					Type: "hangoutsMeet",
				},
			},
		},
	}

	created, err := client.Events.
		Insert("primary", event).
		ConferenceDataVersion(1).
		Do()
	if err != nil {
		return "", err
	}

	return created.HangoutLink, nil
}

func BookExpertSlot(c *gin.Context, slotID uint) error {

	studentUUID := c.GetString("user_uuid")

	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1️⃣ Lock slot row
	var slot models.AvailabilitySlot
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND is_booked = false", slotID).
		First(&slot).Error

	if err != nil {
		tx.Rollback()
		return errors.New("slot not available")
	}

	// 2️⃣ Create Google Calendar service
	calService, err := calendar.NewService(
		context.Background(),
		option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS_JSON")),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 3️⃣ Create Meet link
	meetLink, err := CreateGoogleMeetLink(
		context.Background(),
		calService,
		slot.StartTime,
		slot.EndTime,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 4️⃣ Create session
	session := models.Session{
		StudentUUID: studentUUID,
		ExpertUUID:  slot.ExpertID,
		SlotID:      slot.ID,
		MeetLink:    meetLink,
		StartTime:   slot.StartTime,
		EndTime:     slot.EndTime,
	}

	if err := tx.Create(&session).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 5️⃣ Mark slot as booked
	if err := tx.Model(&models.AvailabilitySlot{}).
		Where("id = ?", slot.ID).
		Update("is_booked", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
