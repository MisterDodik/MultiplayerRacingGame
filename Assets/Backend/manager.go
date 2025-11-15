package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Manager struct {
	clients ClientList

	events map[string]EventHandler
	otps   RetentionMap
	sync.RWMutex
}

func NewManager() *Manager {
	m := &Manager{
		clients: make(ClientList),
		events:  make(map[string]EventHandler),
		otps:    NewRetentionMap(),
	}

	m.initHandlers()

	return m
}

func (m *Manager) initHandlers() {
	m.events[JoinLobby] = JoinLobbyHandler
	m.events[StartGame] = StartGameHandler
	m.events[ChatroomMsg] = ChatMsgFromClientHandler
}

func (m *Manager) parseEvent(e Event, c *Client) error {
	event, ok := m.events[e.Type]
	if !ok {
		return errors.New("unknown event type")
	}
	if err := event(e, c); err != nil {
		log.Printf("error handling event: %v;  err: %v", e, err)
		return err
	}
	return nil
}

type LoginPayload struct {
	Seed     string `json:"seed"`
	Username string `json:"username"`
}

func (m *Manager) login(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("error decoding login request body %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	count := 0
	for c := range m.clients {
		if c.lobby == payload.Seed {
			count++
		}
	}

	if count >= 8 {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("lobby full"))
		return
	}

	otp, err := m.otps.NewOTP(payload.Username, payload.Seed)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	m.otps[otp.key] = *otp

	type otpResponse struct {
		OTP string `json:"otp"`
	}
	res := otpResponse{
		OTP: otp.key,
	}
	data, err := json.Marshal(res)

	if err != nil {
		log.Printf("error marshaling login response %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
	var userName string
	var lobby string

	otp := r.URL.Query().Get("otp")
	if otp == "" || !m.otps.ValidateOTP(otp, &userName, &lobby) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading the conn: ", err)
		return
	}
	log.Println("new connection")

	client := NewClient(conn, m)
	client.username = userName
	client.lobby = lobby

	m.addClient(client)

	go client.ReadMessage()
	go client.WriteMessage()

	lbJson, _ := json.Marshal(lobby)

	evt := Event{
		Type:    JoinLobby,
		Payload: json.RawMessage(lbJson),
	}
	if err := m.parseEvent(evt, client); err != nil {
		m.removeClient(client)
	}
}

func (m *Manager) addClient(client *Client) {
	m.Lock()
	defer m.Unlock()

	m.clients[client] = true

	log.Println("new client")
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[client]; ok {

		client.conn.Close()
		delete(m.clients, client)

		log.Println("client removed")
	}
}
