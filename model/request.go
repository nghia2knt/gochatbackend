package model

type CreateUserForm struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendMessageForm struct {
	ConversationId string `json:"conversationId"`
	Content        string `json:"content"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateConversationForm struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}
