package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Maxito7/hotel_backend/internal/domain"
	"github.com/google/uuid"
)

type chatbotRepository struct {
	db *sql.DB
}

func NewChatbotRepository(db *sql.DB) domain.ChatbotRepository {
	return &chatbotRepository{db: db}
}

func (r *chatbotRepository) SaveConversation(conversation *domain.ConversationHistory) error {
	// Generamos ID si no existe
	if conversation.ID == "" {
		conversation.ID = uuid.New().String()
	}

	messagesJSON, err := json.Marshal(conversation.Messages)
	if err != nil {
		return fmt.Errorf("error marshaling messages: %w", err)
	}

	query := `
		INSERT INTO conversation_history 
		(id, cliente_id, messages, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = r.db.Exec(query,
		conversation.ID,
		conversation.ClienteID,
		messagesJSON,
		conversation.CreatedAt,
		conversation.UpdatedAt,
	)

	return err
}

func (r *chatbotRepository) GetConversation(conversationID string) (*domain.ConversationHistory, error) {
	query := `
		SELECT id, cliente_id, messages, created_at, updated_at 
		FROM conversation_history 
		WHERE id = $1
	`

	var conversation domain.ConversationHistory
	var messagesJSON []byte
	var clienteID sql.NullInt64

	err := r.db.QueryRow(query, conversationID).Scan(
		&conversation.ID,
		&clienteID,
		&messagesJSON,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if clienteID.Valid {
		id := int(clienteID.Int64)
		conversation.ClienteID = &id
	}

	if err := json.Unmarshal(messagesJSON, &conversation.Messages); err != nil {
		return nil, fmt.Errorf("error unmarshaling messages: %w", err)
	}

	return &conversation, nil
}

func (r *chatbotRepository) UpdateConversation(conversation *domain.ConversationHistory) error {
	messagesJSON, err := json.Marshal(conversation.Messages)
	if err != nil {
		return fmt.Errorf("error marshaling messages: %w", err)
	}

	query := `
		UPDATE conversation_history 
		SET messages = $1, updated_at = $2, cliente_id = $3
		WHERE id = $4
	`

	_, err = r.db.Exec(query,
		messagesJSON,
		time.Now(),
		conversation.ClienteID,
		conversation.ID,
	)

	return err
}

func (r *chatbotRepository) SaveMessage(clienteID int, contenido string) error {
	query := `
		INSERT INTO mensaje (contenido, clienteid, fecharegistro) 
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, contenido, clienteID, time.Now())
	return err
}

func (r *chatbotRepository) GetClientConversations(clienteID int) ([]domain.ConversationHistory, error) {
	query := `
		SELECT id, cliente_id, messages, created_at, updated_at 
		FROM conversation_history 
		WHERE cliente_id = $1 
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(query, clienteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []domain.ConversationHistory
	for rows.Next() {
		var conv domain.ConversationHistory
		var messagesJSON []byte
		var clienteID sql.NullInt64

		err := rows.Scan(
			&conv.ID,
			&clienteID,
			&messagesJSON,
			&conv.CreatedAt,
			&conv.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if clienteID.Valid {
			id := int(clienteID.Int64)
			conv.ClienteID = &id
		}

		if err := json.Unmarshal(messagesJSON, &conv.Messages); err != nil {
			return nil, fmt.Errorf("error unmarshaling messages: %w", err)
		}

		conversations = append(conversations, conv)
	}

	return conversations, nil
}
