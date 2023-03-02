package database

type User struct {
	Id          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Key         []byte `json:"key"`
}
