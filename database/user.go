package database

import "fmt"

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Key         []byte `json:"key"`
}

func GetUser(id string) (*User, error) {
	user, ok := Users[id]
	if !ok {
		return nil, fmt.Errorf("no user found with id %s", id)
	} else {
		return &user, nil
	}
}

func GetUserByUsername(username string) (*User, error) {
	for _, user := range Users {
		if user.Username == username {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("no user found with username %s", username)
}
