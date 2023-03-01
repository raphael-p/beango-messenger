package database

import (
	"time"
)

type Session struct {
	Id         string    `json:"id"`
	UserId     string    `json:"userId"`
	ExpiryDate time.Time `json:"expiryDate"`
}

func AddSession(session Session) {
	Sessions[session.Id] = session
}

func GetSession(id string) *Session {
	session, ok := Sessions[id]
	if !ok {
		return nil
	} else {
		return &session
	}
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
