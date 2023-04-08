package database

import (
	"fmt"
	"time"
)

type Session struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userID"`
	ExpiryDate time.Time `json:"expiryDate"`
}

func GetSession(id string) *Session {
	session, ok := Sessions[id]
	if !ok {
		return nil
	} else {
		return &session
	}
}

func SetSession(session Session) {
	if session, err := GetSessionByUserID(session.UserID); err == nil {
		DeleteSession(session.ID)
	}

	Sessions[session.ID] = session
}

func DeleteSession(id string) {
	delete(Sessions, id)
}

func CheckSession(id string) (*Session, bool) {
	if id == "" {
		return nil, false
	}
	session := GetSession(id)
	if session == nil {
		return nil, false
	}
	if session.ExpiryDate.Before(time.Now().UTC()) {
		DeleteSession(session.ID)
		return nil, false
	}
	return session, true
}

func GetSessionByUserID(userID string) (*Session, error) {
	for _, session := range Sessions {
		if session.UserID == userID {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("no session found for user ID %s", userID)
}
