package models

import "time"

type User struct {
	UserId    string
	UserName  string
	FriendsList []string // list of user IDs
	GroupList []string // list of group IDs
}

type Group struct {
	GroupId      string
	UserList     []string //userid list
	ChatHistory []GroupMessage
}

type Message struct {
	MessageId   string
	SenderId    string
	ReceiverId  string
	MessageType string //use short int in future
	Content     string //convert into bit-representation in the future
	TimeStamp   time.Time
}

type GroupMessage struct {
	GroupId		string
	MessageId   string
	SenderId    string
	MessageType string //use short int in future
	Content     string //convert into bit-representation in the future
	TimeStamp   time.Time
}

type UserLoginCredentials struct {
	UserName string	`json:"username"`
	Password string	`json:"password"`
}

type UserSignupCredentials struct {
	UserName string	`json:"username"`
	Password string	`json:"password"`
	Email string `json:"email"`
}