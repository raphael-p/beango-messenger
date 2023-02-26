package database

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
