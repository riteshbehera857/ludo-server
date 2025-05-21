package main

import (
	"flag"
	"log"
	"ludo"
	"messaging/common"
	"messaging/rest"
	"messaging/socket"
	"metagame/gameserver/config"
	"metagame/gameserver/helpers"
	"strings"
	"sync"
)

var gameServiceMap = make(map[string]common.GameService)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config.GetConfig()

	portPtr := flag.Int("port", 4000, "port number for socket socket")
	restServerPortPtr := flag.String("restPort", "4001", "port number for rest server")
	playerId := flag.String("playerId", "player1", "player id for the game")

	flag.Parse()

	playerIds := strings.Split(*playerId, ",")

	for _, id := range playerIds {
		_, err := socket.CreateToken(id)
		if err != nil {
			log.Printf("Error %s when creating token for player %s", err, id)
			continue
		}
		// log.Printf("Token %s : %s", id, token)
	}

	// log.Printf("Starting socket on port %d", *portPtr)

	_, err := helpers.InitializeNewMongoClient(cfg.MongoURI)

	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Add game services to the map
	gameServiceMap["ludo"] = &ludo.LudoGameService{}

	for _, v := range gameServiceMap {
		// log.Printf("Game service %s added", k)
		v.StartBoardManagement()
	}

	var wg sync.WaitGroup
	wg.Add(2)

	// Start Socket Server
	go func() {
		defer wg.Done()
		if err := socket.StartSocketServer(*portPtr, gameServiceMap); err != nil {
			// log.Printf("Socket server error: %v", err)
		}
	}()

	// Start REST API Server
	go func() {
		defer wg.Done()
		if err := rest.StartRESTApiServer(*restServerPortPtr); err != nil {
			log.Printf("REST API server error: %v", err)
		}
	}()

	// Wait for both servers
	wg.Wait()

}
