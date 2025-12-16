package handlers

import (
	"github.com/MisterDodik/MultiplayerGame/internal/events"
	"github.com/MisterDodik/MultiplayerGame/internal/network"
)

func EndGameHandler(e events.Event, c *network.Client) error {

	//a nek stoji ovak funkcija, mzd ce mi trebati za nesto sa client strane, samo u types uncommentuj

	//na client strani ce da:
	//	pokaze end screen ui
	//  izbaci igraca u lobby, kad on pritisne dugme back to lobby

	return nil
}
