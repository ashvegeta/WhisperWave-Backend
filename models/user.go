package models

import "time"

type User struct {
	UserId      string               `json:"userId"`
	UserName    string               `json:"userName"`
	FriendsList []string             `json:"friendsList"`
	GroupList   []string             `json:"groupList"`
	ChatHistory map[string][]Message `json:"userChatHistory"`
}

type Group struct {
	GroupId     string    `json:"groupId"`
	GroupName   string    `json:"groupName"`
	UserList    []string  `json:"userList"`
	ChatHistory []Message `json:"groupChatHistory"`
}

type Message struct {
	MessageId   string    `json:"messageId"`
	SenderId    string    `json:"senderId"`
	GroupId     string    `json:"groupId"`
	ReceiverIds []string  `json:"receiverIds"`
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
