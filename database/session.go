package database

import (
	"net/http"
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

func CheckSession(cookie *http.Cookie) (*Session, bool) {
	if cookie == nil {
		return nil, false
	}
	session := GetSession(cookie.Value)
	if session == nil {
		return nil, false
	}
	if session.ExpiryDate.Before(time.Now()) {
		return session, false
	}
	return session, true
}
