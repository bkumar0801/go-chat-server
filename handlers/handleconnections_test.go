package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chat-server/context"
	"github.com/gorilla/websocket"
)

func TestMakeUserOnline(t *testing.T) {
	ctx := context.NewAppContext()
	msg := context.Message{
		Email:    "xyz@xyz.com",
		Username: "Test",
		Message:  "Login",
	}
	makeUserOnline(ctx, &msg, msg.Username)

	if msg.Message != "Test is online" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Test is online' \n\t\t actual = %s", msg.Message)
	}
	if msg.Username != "Server" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Server' \n\t\t actual = %s", msg.Username)
	}
	if val, _ := ctx.Users["Test"]; val != 1 {
		t.Errorf("Expectation mismatched: \n\t\t expected = 1 \n\t\t actual = %d", val)
	}

	msg = context.Message{
		Email:    "xyz@xyz.com",
		Username: "Test",
		Message:  "Login",
	}
	makeUserOnline(ctx, &msg, msg.Username)
	if msg.Message != "Test is already logged in" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Test is already logged in' \n\t\t actual = %s", msg.Message)
	}
	if msg.Username != "Server" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Server' \n\t\t actual = %s", msg.Username)
	}
	if val, _ := ctx.Users["Test"]; val != 2 {
		t.Errorf("Expectation mismatched: \n\t\t expected = 1 \n\t\t actual = %d", val)
	}

	close(ctx.Broadcast)
}

func TestMakeUserOffline(t *testing.T) {

	ctx := context.NewAppContext()
	msg := context.Message{
		Email:    "xyz@xyz.com",
		Username: "Test",
		Message:  "Login",
	}
	makeUserOnline(ctx, &msg, msg.Username)
	if msg.Message != "Test is online" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Test is online' \n\t\t actual = %s", msg.Message)
	}
	if msg.Username != "Server" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Server' \n\t\t actual = %s", msg.Username)
	}
	if val, _ := ctx.Users["Test"]; val != 1 {
		t.Errorf("Expectation mismatched: \n\t\t expected = 1 \n\t\t actual = %d", val)
	}

	msg = context.Message{
		Email:    "xyz@xyz.com",
		Username: "Test",
		Message:  "Login",
	}
	makeUserOffline(ctx, &msg, msg.Username)

	if msg.Message != "Test is offline" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Test is online' \n\t\t actual = %s", msg.Message)
	}
	if msg.Username != "Server" {
		t.Errorf("Expectation mismatched: \n\t\t expected = 'Server' \n\t\t actual = %s", msg.Username)
	}
	if val, _ := ctx.Users["Test"]; val != 0 {
		t.Errorf("Expectation mismatched: \n\t\t expected = 1 \n\t\t actual = %d", val)
	}

	close(ctx.Broadcast)
}

func TestHandleConnection(t *testing.T) {
	appCtx := context.NewAppContext()
	testServer := httptest.NewServer(http.HandlerFunc(HandleConnections(appCtx)))
	defer testServer.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1/echo
	u := "ws" + strings.TrimPrefix(testServer.URL, "http") + "/echo"

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	message := context.Message{
		Email:    "xyz@xyz.com",
		Username: "Test",
		Message:  "test message",
	}
	data, _ := json.Marshal(message)
	err = ws.WriteMessage(websocket.BinaryMessage, data)
	if err != nil {
		t.Fatalf("%v", err)
	}
	want := context.Message{
		Email:    "xyz@xyz.com",
		Username: "Server",
		Message:  "Test is online",
	}
	got := <-appCtx.Broadcast
	if want != got {
		t.Errorf("\nUnexpected response from server: \n\t\t expected: %q, \n\t\t actual: %q", want, got)
	}
}
