package models

// NOTE: chat history for users and groups are stored differently in the DynamoDB, there is no
// need to add chat history to "User" OR "Group" struct
type User struct {
	UserId      string   `json:"userId" dynamodbav:"ID"`
	EmailID     string   `json:"emailId" dynamodbav:"EmailID"`
	UserName    string   `json:"userName" dynamodbav:"UserName"`
	Password    string   `json:"password" dynamodbav:"Password"`
	FriendsList []string `json:"friendsList" dynamodbav:"FriendsList,stringset"`
	GroupList   []string `json:"groupList" dynamodbav:"GroupsList,stringset"`
}

type Group struct {
	GroupId   string   `json:"groupId" dynamodbav:"ID"`
	GroupName string   `json:"groupName" dynamodbav:"GroupName"`
	UserList  []string `json:"userList" dynamodbav:"UserList,stringset"`
}

type Message struct {
	MessageId   string   `json:"messageId"`
	SenderId    string   `json:"senderId"`
	GroupId     string   `json:"groupId"`
	ReceiverIds []string `json:"receiverIds"`
	MessageType string   `json:"messageType"`
	Content     string   `json:"content"`
	TimeStamp   string   `json:"timeStamp"`
}

// only use this to verify credentials, dont store directly
type UserLoginCredentials struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
}

// only use this to verify credentials, dont store directly
type UserSignupCredentials struct {
	UserId   string `json:"userId"`
	Password string `json:"password"`
	UserName string `json:"username"`
	Email    string `json:"email"`
}
