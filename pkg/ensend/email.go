package ensend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	ProjectID string
	ProjectSecret string
}


const apiURL = "https://api.smtpexpress.com/send"

type SenderInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Payload struct {
	Subject    string      `json:"subject"`
	Message    string      `json:"message"`
	Sender     SenderInfo  `json:"sender"`
	Recipients string      `json:"recipients"`
}

func NewConfig(projectID, projectSecret string) *Config {
	return &Config{
		ProjectID:     projectID,
		ProjectSecret: projectSecret,
	}
}


func NewPayload(subject, message, senderEmail, senderName, recipients string) *Payload {
	return &Payload{
		Subject: subject,
		Message: message,
		Sender: SenderInfo{
			Email: senderEmail,
			Name:  senderName,
		},
		Recipients: recipients,
	}
}

func SendEmail(cfg *Config, payload *Payload) error {
client := &http.Client{
		Timeout: 10 * time.Second, // always have a timeout; the universe is chaotic
	}
	
	// Construct the request payload
	//construct the sender info
	sender := map[string]string{
		"email": payload.Sender.Email,
		"name":  payload.Sender.Name,
	}
	
	//construct the full payload
	requestPayload := map[string]interface{}{
		"subject":    payload.Subject,
		"message":    payload.Message,
		"sender":     sender,
		"recipients": payload.Recipients,
	}
	jsonBody, _ := json.Marshal(requestPayload)

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("X-Project-ID", cfg.ProjectID)
	req.Header.Set("Authorization","Bearer " + cfg.ProjectSecret)

	postResp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer postResp.Body.Close()

	var parsed any
	if err := json.NewDecoder(postResp.Body).Decode(&parsed); err != nil {
		panic(err)
	}
	fmt.Println("Response from SMTP Express:", parsed)
	return nil
}

