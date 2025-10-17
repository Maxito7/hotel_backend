package domain

import "time"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Message        string       `json:"message"`
	ConversationID *string      `json:"conversationId,omitempty"`
	ClienteID      *int         `json:"clienteId,omitempty"`
	Context        *ChatContext `json:"context,omitempty"`
	// UseWeb: nil = auto (service decides), true = force web search, false = disable web search
	UseWeb *bool `json:"useWeb,omitempty"`
}

type ChatContext struct {
	FechaEntrada    *string `json:"fechaEntrada,omitempty"`
	FechaSalida     *string `json:"fechaSalida,omitempty"`
	CantidadAdultos *int    `json:"cantidadAdultos,omitempty"`
	CantidadNinhos  *int    `json:"cantidadNinhos,omitempty"`
}

type ChatResponse struct {
	Message          string                 `json:"message"`
	ConversationID   string                 `json:"conversationId"`
	SuggestedActions []string               `json:"suggestedActions,omitempty"`
	RequiresHuman    bool                   `json:"requiresHuman"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

type ConversationHistory struct {
	ID        string        `json:"id"`
	ClienteID *int          `json:"clienteId,omitempty"`
	Messages  []ChatMessage `json:"messages"`
	CreatedAt time.Time     `json:"createdAt"`
	UpdatedAt time.Time     `json:"updatedAt"`
}

type ChatbotRepository interface {
	SaveConversation(conversation *ConversationHistory) error
	GetConversation(conversationID string) (*ConversationHistory, error)
	UpdateConversation(conversation *ConversationHistory) error
	SaveMessage(clienteID int, contenido string) error
	GetClientConversations(clienteID int) ([]ConversationHistory, error)
}
