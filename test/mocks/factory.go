package mocks

import (
	"time"

	"github.com/google/uuid"
	"github.com/raphael-p/beango/database"
)

const (
	ADMIN_ID       int64 = 1
	ADMIN_USERNAME       = "the_admin"
	PASSWORD             = "123abc*"
	HASH                 = "$2y$04$8QoTLjUMGtnr4lNeA0DtduhEshmvbDbmEzW/G9IkV/9mr576xX//K"
)

func MakeUser() *database.User {
	return &database.User{
		Username:    "john.doe.69",
		DisplayName: "Johnny D",
		Key:         []byte(HASH),
	}
}

func MakeUser2() *database.User {
	return &database.User{
		Username:    "miltonb",
		DisplayName: "Blake Milton",
		Key:         []byte(HASH),
	}
}

func MakeAdminUser() *database.User {
	return &database.User{
		ID:          ADMIN_ID,
		Username:    ADMIN_USERNAME,
		DisplayName: "Administrator",
		Key:         []byte(HASH),
	}
}

func MakeSession(userID int64) database.Session {
	return database.Session{
		ID:         uuid.NewString(),
		UserID:     userID,
		ExpiryDate: time.Now().UTC().Add(time.Hour),
	}
}

func MakePrivateChat() *database.Chat {
	return &database.Chat{
		Type: database.PRIVATE_CHAT,
	}
}

func MakeMessage(userID, chatID int64) *database.MessageDatabase {
	return &database.MessageDatabase{
		UserID:  userID,
		ChatID:  chatID,
		Content: "Lorem Ipsum Dolor",
	}
}
