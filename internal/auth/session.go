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
	data.ExpireOn = 86400 // 24 hour Expire Time

	sessionId, err := repository.SaveSession(data)

	if err != nil {
		return "", err
	}

	return sessionId, nil

}

func SignOut(sessionId string) (string, error) {
	id, err := repository.DeleteSession(sessionId)

	if err != nil {
		return "", err
	}

	return id, nil

}

func GetSession(session string) (*models.UserSession, error) {

	res, err := repository.FindSession(session)

	if err != nil {
		return nil, err
	}

	return res, nil

}
