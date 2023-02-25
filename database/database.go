package database

var Users = make(map[string]User)
var Chats = make(map[string]Chat)
var Messages = make(map[string]Message)

func GetUserByUsername(username string) (User, bool) {
	for _, user := range Users {
		if user.Username == username {
			return user, true
		}
	}
	return User{}, false
}
