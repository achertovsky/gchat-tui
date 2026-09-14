package googlechat

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListSpacesPaginates(t *testing.T) {
	var requests int
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requests++
		if request.URL.Path != "/v1/spaces" {
			t.Errorf("path = %q, want /v1/spaces", request.URL.Path)
		}
		if request.URL.Query().Get("pageToken") == "" {
			writeJSON(t, writer, map[string]any{
				"spaces":        []map[string]string{{"name": "spaces/one", "displayName": "One", "spaceType": "SPACE"}},
				"nextPageToken": "second-page",
			})
			return
		}
		writeJSON(t, writer, map[string]any{
			"spaces": []map[string]string{{"name": "spaces/two", "displayName": "Two", "spaceType": "DIRECT_MESSAGE"}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	spaces, err := client.ListSpaces(context.Background())
	if err != nil {
		t.Fatalf("list spaces: %v", err)
	}
	if requests != 2 || len(spaces) != 2 {
		t.Fatalf("requests = %d, spaces = %#v; want two pages and spaces", requests, spaces)
	}
	if spaces[1].Type != "DIRECT_MESSAGE" {
		t.Errorf("space type = %q, want DIRECT_MESSAGE", spaces[1].Type)
	}
}

func TestListMessagesPaginates(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.URL.Path != "/v1/spaces/example/messages" {
			t.Errorf("path = %q, want message endpoint", request.URL.Path)
		}
		if request.URL.Query().Get("pageToken") == "" {
			writeJSON(t, writer, map[string]any{
				"messages":      []map[string]string{{"name": "spaces/example/messages/one", "text": "first"}},
				"nextPageToken": "second-page",
			})
			return
		}
		writeJSON(t, writer, map[string]any{
			"messages": []map[string]string{{"name": "spaces/example/messages/two", "text": "second"}},
		})
	}))
	defer server.Close()

	client := newTestClient(t, server)
	messages, err := client.ListMessages(context.Background(), "spaces/example")
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 2 || messages[0].Text != "second" {
		t.Errorf("messages = %#v, want two paginated messages", messages)
	}
}

func TestCreateMessageAndSanitizeAPIErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodPost {
			var message map[string]string
			if err := json.NewDecoder(request.Body).Decode(&message); err != nil {
				t.Fatal(err)
			}
			if message["text"] != "hello" {
				t.Errorf("text = %q, want hello", message["text"])
			}
			writeJSON(t, writer, map[string]string{"name": "spaces/example/messages/new", "text": "hello"})
			return
		}
		http.Error(writer, `{"error":{"code":403,"message":"credential data must not be exposed"}}`, http.StatusForbidden)
	}))
	defer server.Close()

	client := newTestClient(t, server)
	message, err := client.CreateMessage(context.Background(), "spaces/example", "hello")
	if err != nil {
		t.Fatalf("create message: %v", err)
	}
	if message.Name != "spaces/example/messages/new" {
		t.Errorf("message = %#v", message)
	}

	_, err = client.ListSpaces(context.Background())
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %v, want APIError", err)
	}
	if apiErr.StatusCode != http.StatusForbidden || apiErr.Error() == "credential data must not be exposed" {
		t.Errorf("API error = %#v", apiErr)
	}
}

func newTestClient(t *testing.T, server *httptest.Server) *Client {
	t.Helper()
	client, err := New(server.Client())
	if err != nil {
		t.Fatal(err)
	}
	client.service.BasePath = server.URL + "/"
	return client
}

func writeJSON(t *testing.T, writer http.ResponseWriter, value any) {
	t.Helper()
	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(value); err != nil {
		t.Fatal(err)
	}
}
