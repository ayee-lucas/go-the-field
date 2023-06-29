package auth

import (
	"time"

	"github.com/alopez-2018459/go-the-field/internal/models"
	"github.com/alopez-2018459/go-the-field/internal/repository"
	"github.com/google/uuid"
)

func GenerateSession(data *models.UserSession) (string, error) {

	newSessionId := uuid.NewString()

	data.ID = newSessionId
	data.CreatedAt = time.Now()
	data.ExpireOn = time.Now().Add(time.Hour * 24).Unix() // 24 hour Expire Time

	sessionId, err := repository.SaveSession(data)

	if err != nil {
		return "", err
	}

	return sessionId, nil

}
