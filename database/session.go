package database

import (
	"fmt"
	"time"
)

type Session struct {
	ID         string    `json:"id"`
	UserID     int       `json:"userID"`
	ExpiryDate time.Time `json:"expiryDate"`
}

func (conn *MongoConnection) GetSession(id string) *Session {
	session, ok := Sessions[id]
	if !ok {
		return nil
	} else {
		return &session
	}
}

func (conn *MongoConnection) SetSession(session Session) {
	if session, err := conn.GetSessionByUserID(session.UserID); err == nil {
		conn.DeleteSession(session.ID)
	}

	Sessions[session.ID] = session
}

func (conn *MongoConnection) DeleteSession(id string) {
	delete(Sessions, id)
}

func (conn *MongoConnection) CheckSession(id string) (*Session, bool) {
	if id == "" {
		return nil, false
	}
	session := conn.GetSession(id)
	if session == nil {
		return nil, false
	}
	if session.ExpiryDate.Before(time.Now().UTC()) {
		conn.DeleteSession(session.ID)
		return nil, false
	}
	return session, true
}

func (conn *MongoConnection) GetSessionByUserID(userID int) (*Session, error) {
	for _, session := range Sessions {
		if session.UserID == userID {
			return &session, nil
		}
	}
	return nil, fmt.Errorf("no session found for user ID %d", userID)
}
