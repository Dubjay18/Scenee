package services

import (
	"sync"

	"github.com/Dubjay18/scenee/pkg/ensend"
)







type NotificationConfig struct {
	// Configuration fields for notification service
EnSendProjectID     string
	EnSendProjectSecret string
}

type NotificationService struct {
	// Add dependencies here, e.g., email sender, push notification service, etc.
	cfg *ensend.Config
}

func NewNotificationService(cfg NotificationConfig) *NotificationService {
	return &NotificationService{
		cfg: ensend.NewConfig(
			cfg.EnSendProjectID,
			cfg.EnSendProjectSecret,
		),
	}
}

type INotificationService interface {
	// Define methods for sending notifications
	SendEmailNotification(to string, subject string, body string) error
	SendPushNotification(deviceToken string, title string, message string) error
}
func (s *NotificationService) SendEmailNotification(to string, subject string, body string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		
		si := ensend.SenderInfo{
			Email: "scene-a8f3af@ensend.me",
			Name:  "Scenee Support",
		}
		cfg := ensend.NewConfig(s.cfg.ProjectID, s.cfg.ProjectSecret)
		payload := ensend.NewPayload(subject, body, si.Email, si.Name, to)
		
		if err := ensend.SendEmail(cfg, payload); err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	// Wait for the goroutine to complete
	wg.Wait()
	close(errChan)

	// Check if there was an error
	if err := <-errChan; err != nil {
		return err
	}

	return nil
}