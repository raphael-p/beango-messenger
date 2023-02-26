package database

import (
	"fmt"
)

var Users = make(map[string]User)
var Chats = make(map[string]Chat)
var Messages = make(map[string]Message)
var Sessions = make(map[string]Session)

func GetUserByUsername(username string) *User {
	for _, user := range Users {
		if user.Username == username {
			return &user
		}
	}
	return nil
}

func GetUser(id string) (*User, error) {
	user, ok := Users[id]
	if !ok {
		return nil, fmt.Errorf("no user found with id %s", id)
	} else {
		return &user, nil
	}
}
