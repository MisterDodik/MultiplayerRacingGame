package network

type GameServerList map[string]*GameServer

type GameServer struct {
	Clients ClientList
}

func (m *Manager) NewGameServer(lobbyName string) *GameServer {
	gs := &GameServer{
		Clients: make(ClientList),
	}
	m.Games[lobbyName] = gs
	return gs
}

/*
	client salje event u formatu:
	{
		type: updatePos
		payload: {
					input vector2
				}
	}

	onda se ovdje pozove funkcija koja for loopom prodje kroz sve igrace u lobiju
	calculate_pos(event.payload.input)

	broadcasttoeveryone(client.newpos)


	func calculate_pos(input vector2){
		newPosX = client.currentX + input.x * speed
		newPosY = client.currentY + input.y * speed


		//cekira collision
		for p in players_in_lobby{
			if p = client
				continue
			if rastojanje(p.posX, newPosX) < sizeof_client + sizeof_p{
				newPosX = client.currentX
			}
			if rastojanje(p.posY, newPosY) < sizeof_client + sizeof_p{
				newPosY = client.currentY
			}
		}

		return (newPosX, newPosY)
	}
*/

func (c *Client) UpdatePlayer() {
}
