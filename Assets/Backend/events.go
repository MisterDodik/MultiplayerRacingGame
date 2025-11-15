package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"unicode"
)

type EventHandler func(Event, *Client) error
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

var (
	JoinLobby          = "join_lobby"
	PopulateLobby      = "populate_lobby"
	StartGame          = "start_game"
	ChatroomMsg        = "chatroom_msg"
	BroadcastToClients = "server_msg"
)

func JoinLobbyHandler(e Event, c *Client) error {
	var lobbyName string
	err := json.Unmarshal(e.Payload, &lobbyName)
	if err != nil {
		return err
	}

	if len(lobbyName) < 1 {
		return fmt.Errorf("lobby name too short")
	}
	for _, r := range lobbyName {
		if !unicode.IsDigit(r) {
			return fmt.Errorf("only numbers are allowed in lobby name, but found %q", r)
		}
	}

	c.lobby = lobbyName

	jsonMsg, err := json.Marshal(lobbyName)
	if err != nil {
		return err
	}

	responseEvt := Event{
		Type:    JoinLobby,
		Payload: jsonMsg,
	}
	responseEvt.broadcastMessageToSingleClient(c)

	if err := populateLobby(c); err != nil {
		return err
	}
	return nil
}

func populateLobby(c *Client) error {
	type LobbyPlayer struct {
		Username string `json:"username"`
	}

	data := LobbyPlayer{
		Username: c.username,
	}
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Println("error parsing populate lobby data")
		return err
	}

	evt := Event{
		Type:    PopulateLobby,
		Payload: jsonData,
	}
	evt.broadcastMessageToAllClients(c)

	return nil
}

func StartGameHandler(e Event, c *Client) error {
	log.Println("game started")
	return nil
}

type ChatMessageToClient struct {
	From    string `json:"from"`
	SentAt  string `json:"sentAt"`
	Message string `json:"message"`
}

func ChatMsgFromClientHandler(e Event, c *Client) error {
	var clientMsg string
	err := json.Unmarshal(e.Payload, &clientMsg)
	if err != nil {
		return err
	}

	if len(clientMsg) < 1 {
		return fmt.Errorf("message too short")
	}
	log.Println(clientMsg)

	response := ChatMessageToClient{
		From:    c.username,
		Message: clientMsg,
		SentAt:  time.Now().Format(time.TimeOnly),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	responseEvt := Event{
		Type:    ChatroomMsg,
		Payload: jsonResponse,
	}
	responseEvt.broadcastMessageToAllClients(c)
	return nil
}

func (e Event) broadcastMessageToAllClients(c *Client) {
	sendEventJSON, err := json.Marshal(&e)
	if err != nil {
		log.Println("error marshaling new send event ", err)
		return
	}

	log.Println(string(sendEventJSON))
	for client := range c.manager.clients {
		if client.lobby == c.lobby {
			client.egress <- sendEventJSON
		}
	}
}
func (e Event) broadcastMessageToSingleClient(c *Client) {
	sendEventJSON, err := json.Marshal(&e)
	if err != nil {
		log.Println("error marshaling new send event ", err)
		return
	}
	c.egress <- sendEventJSON
}
