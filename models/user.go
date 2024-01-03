package models

import "time"

type User struct {
	UserId      string   `json:"userId"`
	UserName    string   `json:"userIName"`
	FriendsList []string `json:"friendsList"`
	GroupList   []string `json:"groupList"`
}

type Group struct {
	GroupId     string         `json:"groupId"`
	UserList    []string       `json:"userList"`
	ChatHistory []GroupMessage `json:"chatHistory"`
}

type Message struct {
	MessageId   string    `json:"messageId"`
	SenderId    string    `json:"senderId"`
	ReceiverId  string    `json:"receiverId"`
	MessageType string    `json:"messageType"`
	Content     string    `json:"content"`
	TimeStamp   time.Time `json:"timeStamp"`
}

type GroupMessage struct {
	GroupId     string    `json:"groupId"`
	MessageId   string    `json:"messageId"`
	SenderId    string    `json:"senderId"`
	MessageType string    `json:"messageType"`
	Content     string    `json:"content"`
	TimeStamp   time.Time `json:"timeStamp"`
}

type UserLoginCredentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type UserSignupCredentials struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
