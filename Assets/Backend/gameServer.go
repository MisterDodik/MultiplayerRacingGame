package main

type GameServerList map[string]*GameServer

type GameServer struct {
	lobby string
	clients ClientList
}

func NewGameServer() *GameServer {
	return &GameServer{
		players: make(ClientList),
	}
}
