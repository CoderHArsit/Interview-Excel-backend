package controllers

import (
	"context"
	"errors"
	"interviewexcel-backend-go/config"
	"interviewexcel-backend-go/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/api/calendar/v3"
	"gorm.io/gorm/clause"
	logger "interviewexcel-backend-go/pkg/errors"
)

func CreateGoogleMeetLink(
	ctx context.Context,
	srv *calendar.Service,
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
			},
		},
	}

	createdEvent, err := srv.Events.
		Insert("primary", event).
		ConferenceDataVersion(1).
		Do()

	if err != nil {
		logger.Error("error in creating google meet link: ", err)
		return "", err
	}

	logger.Infof("event id=%s hangout=%s conference=%+v",
		createdEvent.Id,
		createdEvent.HangoutLink,
		createdEvent.ConferenceData,
	)
	return createdEvent.HangoutLink, nil
}

func BookExpertSlot(c *gin.Context, slotID uint) error {

	var (
		tx                   = config.DB.Begin()
		sessionRepo          = models.InitSessionRepo(tx)
		AvailabilitySlotRepo = models.InitAvailabilitySlotRepo(tx)
	)
	studentUUID := c.GetString("user_uuid")

	if tx.Error != nil {
		logger.Error("error in starting transaction: ", tx.Error)
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var slot models.AvailabilitySlot
	err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND status = ?", slotID, models.SlotAvailable).
		First(&slot).Error

	if err != nil {
		tx.Rollback()
		logger.Error("slot not available: ", err)
		return errors.New("slot not available")
	}

	//FIXME: Temporarily disabling Google Meet link creation
	//FEATURE: ADDING PROVIDERS FOR MEET LINK CREATION

	// // 2️⃣ Create Google Calendar service
	// calService, err := calendar.NewService(
	// 	context.Background(),
	// 	option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS_JSON")),
	// )
	// if err != nil {
	// 	tx.Rollback()
	// 	logger.Error("error in creating calendar service: ", err)
	// 	return err
	// }

	// 3️⃣ Create Meet link
	// meetLink, err := CreateGoogleMeetLink(
	// 	context.Background(),
	// 	calService,
	// 	slot.StartTime,
	// 	slot.EndTime,
	// )
	// if err != nil {
	// 	tx.Rollback()
	// 	logger.Error("error in creating google meet link: ", err)
	// 	return err
	// }

	session := &models.Session{
		SessionUUID: uuid.New().String(),
		ExpertUUID:  slot.ExpertID,
		StudentUUID: studentUUID,
		SlotID:      slot.ID,
		StartTime:   slot.StartTime,
		EndTime:     slot.EndTime,
		Status:      "scheduled",
	}

	session.MeetLink = generateJitsiMeetLink(
		session.ExpertUUID,
		session.StudentUUID,
		session.StartTime,
	)

	if err := sessionRepo.Create(session); err != nil {
		logger.Error("error in creating session: ", err)
		tx.Rollback()
		return err
	}

	// 5️⃣ Mark slot as booked
	err = AvailabilitySlotRepo.UpdateWithTx(
		tx,
		&models.AvailabilitySlot{
			Status: string(models.SlotBooked),
		}, &models.AvailabilitySlot{
			ID: slot.ID,
		})
	if err != nil {
		logger.Error("error in marking slot as booked: ", err)
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
