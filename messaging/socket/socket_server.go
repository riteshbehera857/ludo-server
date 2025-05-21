package socket

import (
	"log"
	"messaging/common"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	readWait     = 15 * time.Second // Max time to wait for client response
	pingInterval = 5 * time.Second  // Frequency of server pings
	writeTimeout = 5 * time.Second  // Time to wait for ping/pong writes
)

var playerConnections = make(map[string][]*Connection)
var gameServiceMap = make(map[string]common.GameService)
var boardPlayerMap = make(map[string][]string)

var gameService common.GameService

type webSocketHandler struct {
	upgrader websocket.Upgrader
}

func writeResponse(msg string, c *websocket.Conn) bool {
	var err = c.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("Error %s when sending message to client", err)
		return true
	}
	return false
}

func (wsh webSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	jwtToken := r.Header.Get("Authorization")

	walletAddress := r.Header.Get("walletAddress")

	// Get query param called boardId from url
	boardId := r.URL.Query().Get("boardId")
	// TODO: Get name from jwt token
	game := r.URL.Query().Get("game")

	var gameService = gameServiceMap[game]

	if gameService == nil {
		// Hack: Default to ludo
		gameService = gameServiceMap["ludo"]
	}

	if jwtToken == "" || walletAddress == "" {
		// log.Println("No token provided")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No jwt token or wallet address is provided"))
		return
	}

	// log.Printf("Token: %s", jwtToken)

	playerId, name, err := VerifyToken(jwtToken)

	if err != nil {
		log.Printf("Error %s when verifying token", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("No token provided"))
		return
	}

	c, err := wsh.upgrader.Upgrade(w, r, nil)

	if err != nil {
		// log.Printf("error %s when upgrading connection to websocket", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	connection := NewConnection(c)

	playerConnections[playerId] = append(playerConnections[playerId], connection)

	// Check if player is already in the board's player list
	playerExists := false

	for _, pid := range boardPlayerMap[boardId] {
		if pid == playerId {
			playerExists = true
			log.Printf("[@ServeHTTP] Player %s already exists in board %s", playerId, boardId)
			break
		}
	}

	if !playerExists {
		boardPlayerMap[boardId] = append(boardPlayerMap[boardId], playerId)
		log.Printf("[@ServeHTTP] Added player %s to board %s", playerId, boardId)
	}

	log.Printf("[@ServeHTTP] New connection established - Player: %s, Board: %s, Game: %s", playerId, boardId, game)
	// log.Printf("[@ServeHTTP] Active players in board %s: %v", boardId, boardPlayerMap[boardId])
	// numConnections := len(playerConnections)
	// log.Printf("[@ServeHTTP] Total active connections: %d", numConnections)

	// TODO: Get name from the platform
	addPlayerError := gameService.AddPlayer(boardId, playerId, name, walletAddress)

	// Configure ping/pong handlers properly
	c.SetPongHandler(func(appData string) error {
		c.SetReadDeadline(time.Now().Add(readWait)) // Use configured readWait
		return nil
	})

	// Set initial read deadline immediately after upgrade
	c.SetReadDeadline(time.Now().Add(readWait))

	go func() {
		ticker := time.NewTicker(pingInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.SetWriteDeadline(time.Now().Add(writeTimeout))
				if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Ping failed for %s: %v", playerId, err)
					c.Close() // This will trigger the read error
					return
				}
			}
		}
	}()

	if addPlayerError != nil {
		log.Printf("Error %s when adding player to game", addPlayerError)
		SendErrorMessage(addPlayerError.Error(), c)
		w.WriteHeader(http.StatusInternalServerError)
		// delete(playerConnections, playerId)
		// delete(boardPlayerMap, boardId)
		c.Close()
		return
	}

	defer func() {
		handleDisconnection(c, boardId, playerId)
	}()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Connection closed unexpectedly: %v", err) // Exit loop on any read error
			}

			// log.Printf("Error %s when reading message from client", err)
			SendErrorMessage("Error %s when reading message from client", c)
			return
		}

		c.SetReadDeadline(time.Now().Add(readWait))

		if mt == websocket.BinaryMessage {
			SendErrorMessage("socket doesn't support binary messages", c)
			err = c.WriteMessage(websocket.TextMessage, []byte("socket doesn't support binary messages"))
			if err != nil {
				// log.Printf("Error %s when sending message to client", err)
			}
			return
		}

		textMessage := string(msg)

		req := common.SocketMessage{}
		parsedReq, err := req.ToObject(textMessage)
		if err != nil {
			SendErrorMessage("Invalid message", c)
			continue
		}
		req = parsedReq.(common.SocketMessage)
		// log.Printf("Receive message %s - %s", req.GetEventName(), string(msg))
		processError := gameService.ProcessMessage(boardId, playerId, req, msg)
		if processError != nil {
			// log.Printf("Error %s when processing message", processError)
			SendErrorMessage(processError.Error(), c)
		}
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/index.html")
}

func StartSocketServer(port int, gsMap map[string]common.GameService) error {

	// ctx, cancel := context.WithCancel(context.Background())

	// defer cancel()

	webSocketHandler := webSocketHandler{
		upgrader: websocket.Upgrader{},
	}
	gameServiceMap = gsMap

	http.Handle("/ws", webSocketHandler)
	http.HandleFunc("/", handleHome)

	log.Printf("Starting socket socket on port %d...", port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	return nil
}

func SendMessage(playerId string, msg common.Message, boardId string) {
	if playerConnections[playerId] == nil {
		// log.Printf("Player %s is not connected", playerId)
		return
	}
	body, err := msg.ToJSON()
	// Put a log here
	// log.Println("Sending message...", body)

	if err != nil {
		// log.Printf("Error %s when marshalling message", err)
		return
	}
	if connections := playerConnections[playerId]; len(connections) > 0 {
		writeResponse(body, connections[len(connections)-1].Conn)
		log.Printf("[%s] Sent Message: [%s] to Player [%s]", boardId, body, playerId)
	}
}

func BroadcastMessage(msg common.Message, boardId string) {
	body, err := msg.ToJSON()

	// Put a log here
	// log.Println("Broadcasting message...", body)

	if err != nil {
		log.Printf("Error %s when marshalling message", err)
		return
	}

	playerIds := boardPlayerMap[boardId]

	if playerIds == nil {
		// log.Printf("No players in board %s", boardId)
		return
	}

	// Send to specified players
	for _, playerId := range playerIds {
		if _, ok := playerConnections[playerId]; ok {
			if connections := playerConnections[playerId]; len(connections) > 0 {
				writeResponse(body, connections[len(connections)-1].Conn)
				log.Printf("[%s] Sent Message: [%s] to Player [%s]", boardId, body, playerId)
			}
		}
	}
}

func handleDisconnection(c *websocket.Conn, boardId string, playerId string) {
	log.Print("Handling disconnection...")
	if c == nil {
		log.Println("Connection is nil")
		return
	}

	// Get the correct game service instance
	gameService := gameServiceMap["ludo"] // Or pass it through parameters
	if gameService == nil {

		log.Println("Game service is nil ")
		gameService = gameServiceMap["ludo"] // Default fallback
	}

	connections := playerConnections[playerId]
	if len(connections) == 0 {
		log.Printf("No connections found for player %s", playerId)
		return
	}

	// Find and remove the specific connection
	newConnections := make([]*Connection, 0)
	for _, conn := range connections {
		if conn.Conn != c {
			newConnections = append(newConnections, conn)
		}
	}

	if len(newConnections) == 0 {
		// Last connection removed
		gameService.HandleDisconnection(boardId, playerId)
		delete(playerConnections, playerId)
		// Remove from board player map
		boardPlayerMap[boardId] = removeFromSlice(boardPlayerMap[boardId], playerId)
	} else {
		playerConnections[playerId] = newConnections
	}

	c.Close()
}

// Helper function to remove from slice
func removeFromSlice(slice []string, item string) []string {
	newSlice := make([]string, 0)
	for _, i := range slice {
		if i != item {
			newSlice = append(newSlice, i)
		}
	}
	return newSlice
}

func SendErrorMessage(errorMessage string, c *websocket.Conn) {
	errMsg, _ := common.NewSocketMessage("error", 500, errorMessage).ToJSON()
	writeResponse(errMsg, c)
}
