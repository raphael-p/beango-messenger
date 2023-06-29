package database

import (
	"fmt"
	"time"
)

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	DisplayName   string    `json:"displayName"`
	Key           []byte    `json:"key"`
	CreatedAt     time.Time `json:"createdAt"`
	LastUpdatedAt time.Time `json:"LastUpdatedAt"`
}

func (conn *MongoConnection) GetUser(id int) (*User, error) {
	user, ok := Users[id]
	if !ok {
		return nil, fmt.Errorf("no user found with id %d", id)
	} else {
		return &user, nil
	}
}

func (conn *MongoConnection) GetUserByUsername(username string) (*User, error) {
	for _, user := range Users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("no user found with username %s", username)
}

func (conn *MongoConnection) SetUser(user *User) *User {
	user.ID = len(Users) + 1
	Users[user.ID] = *user
	return user
}
