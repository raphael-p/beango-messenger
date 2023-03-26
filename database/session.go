package database

import (
	"fmt"
	"time"
)

type Session struct {
	Id         string    `json:"id"`
	UserId     string    `json:"userId"`
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
	if session, err := GetSessionByUserId(session.UserId); err == nil {
		DeleteSession(session.Id)
	}

	Sessions[session.Id] = session
}

func DeleteSession(id string) {
	delete(Sessions, id)
}

func CheckSession(id string) (*Session, bool) {
	session := GetSession(id)
	if session == nil {
		return nil, false
	}
	if session.ExpiryDate.Before(time.Now()) {
		DeleteSession(session.Id)
		return nil, false
	}
	return session, true
}

func GetSessionByUserId(userId string) (*Session, error) {
	for _, session := range Sessions {
		if session.UserId == userId {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("no session found for user ID %s", userId)
}
