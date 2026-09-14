// Package googlechat wraps the Google Chat REST API.
package googlechat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	chatapi "google.golang.org/api/chat/v1"
	"google.golang.org/api/googleapi"
)

const requestTimeout = 30 * time.Second

// Conversation represents a Google Chat space or direct message.
type Conversation struct {
	Name        string
	DisplayName string
	Type        string
}

// Space is a Google Chat conversation resource.
type Space = Conversation

// User represents a Google Chat user.
type User struct {
	Name        string
	DisplayName string
	Type        string
}

// Message represents a text message in a conversation.
type Message struct {
	Name       string
	Sender     User
	Text       string
	CreateTime time.Time
}

// APIError is a sanitized Google Chat API failure.
type APIError struct {
	Operation  string
	StatusCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("Google Chat API %s failed (HTTP %d)", e.Operation, e.StatusCode)
}

// Client provides typed, context-aware access to Google Chat.
type Client struct {
	service *chatapi.Service
}

// New creates a client using an authenticated HTTP client.
func New(httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		return nil, errors.New("authenticated HTTP client is required")
	}

	client := *httpClient
	if client.Timeout == 0 {
		client.Timeout = requestTimeout
	}
	service, err := chatapi.New(&client)
	if err != nil {
		return nil, fmt.Errorf("create Google Chat service: %w", err)
	}
	return &Client{service: service}, nil
}

// ListSpaces returns every space visible to the authenticated user.
func (c *Client) ListSpaces(ctx context.Context) ([]Space, error) {
	var spaces []Space
	var pageToken string
	for {
		response, err := c.service.Spaces.List().
			PageSize(1000).
			PageToken(pageToken).
			Context(ctx).
			Do()
		if err != nil {
			return nil, apiError("list spaces", err)
		}
		for _, space := range response.Spaces {
			spaces = append(spaces, toSpace(space))
		}
		if response.NextPageToken == "" {
			return spaces, nil
		}
		pageToken = response.NextPageToken
	}
}

// ListMessages returns every message in a space.
func (c *Client) ListMessages(ctx context.Context, spaceName string) ([]Message, error) {
	if !validSpaceName(spaceName) {
		return nil, errors.New("space name must use the format spaces/{space}")
	}

	var messages []Message
	var pageToken string
	for {
		response, err := c.service.Spaces.Messages.List(spaceName).
			PageSize(1000).
			PageToken(pageToken).
			Context(ctx).
			Do()
		if err != nil {
			return nil, apiError("list messages", err)
		}
		for _, message := range response.Messages {
			messages = append(messages, toMessage(message))
		}
		if response.NextPageToken == "" {
			return messages, nil
		}
		pageToken = response.NextPageToken
	}
}

// CreateMessage sends a plain-text message to a space.
func (c *Client) CreateMessage(ctx context.Context, spaceName, text string) (Message, error) {
	if !validSpaceName(spaceName) {
		return Message{}, errors.New("space name must use the format spaces/{space}")
	}
	if strings.TrimSpace(text) == "" {
		return Message{}, errors.New("message text must not be empty")
	}

	response, err := c.service.Spaces.Messages.Create(spaceName, &chatapi.Message{Text: text}).
		Context(ctx).
		Do()
	if err != nil {
		return Message{}, apiError("create message", err)
	}
	return toMessage(response), nil
}

func validSpaceName(name string) bool {
	return strings.HasPrefix(name, "spaces/") && len(strings.TrimPrefix(name, "spaces/")) > 0 && !strings.Contains(strings.TrimPrefix(name, "spaces/"), "/")
}

func toSpace(space *chatapi.Space) Space {
	return Space{Name: space.Name, DisplayName: space.DisplayName, Type: space.SpaceType}
}

func toMessage(message *chatapi.Message) Message {
	converted := Message{
		Name:   message.Name,
		Text:   message.Text,
		Sender: toUser(message.Sender),
	}
	if timestamp, err := time.Parse(time.RFC3339, message.CreateTime); err == nil {
		converted.CreateTime = timestamp
	}
	return converted
}

func toUser(user *chatapi.User) User {
	if user == nil {
		return User{}
	}
	return User{Name: user.Name, DisplayName: user.DisplayName, Type: user.Type}
}

func apiError(operation string, err error) error {
	var googleError *googleapi.Error
	if errors.As(err, &googleError) {
		return &APIError{Operation: operation, StatusCode: googleError.Code}
	}
	return fmt.Errorf("Google Chat API %s: %w", operation, err)
}
