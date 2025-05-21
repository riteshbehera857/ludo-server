package board

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"ludo/dice"
	"ludo/ludo_board_constants"
	"ludo/pawn"
	"ludo/player"
	"ludo/quadrant"
	"messaging/common"
	"messaging/socket"
	"metagame/gameserver/config"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Board represents the game board containing quadrants and manages the game state
type Board struct {
	id                         string                           // Unique identifier for the game
	quadrants                  []*quadrant.Quadrant             // Four quadrants (red, green, yellow, blue)
	players                    []*player.Player                 // List of players in the game
	currentTurn                string                           // quadrant.Quadrant which has the current turn
	safePositions              []int                            // List of safe positions on the board
	nextTurn                   string                           // quadrant.Quadrant which will play next
	status                     ludo_board_constants.BoardStatus // Current game status
	autoPlay                   bool
	autoPlayTimer              int
	ticketAmount               int
	rakeAmount                 int
	rakeAmountType             ludo_board_constants.RakeAmountType
	playersRequiredToStartGame int
	expectedMessage            *ExpectedMessage
	diceRolledValue            int
}

type ExpectedMessage struct {
	EventName string
	Quadrant  string
	PlayerId  string
	TStamp    time.Time
	Timeout   time.Duration
	Steps     int
}

// NewBoard creates and initializes a new game board with the given players
// Parameters:
//   - players ([]*player.Player): List of players participating in the game
//
// Returns:
//   - *Board: Pointer to the newly created Board
func NewBoard(boardId string, playersRequiredToStartGame int, autoPlay bool, ticketAmount int, rakeAmount int, rakeAmountType ludo_board_constants.RakeAmountType, autoPlayTimer int) *Board {
	quadrantConfigResult, err := quadrant.InitializeQuadrantConfigIfNotExists(ludo_board_constants.QuadrantsNames, ludo_board_constants.QuadrantsPaths, ludo_board_constants.QuadrantsColors, ludo_board_constants.SafePositions)

	if err != nil {
		fmt.Print("Error: ", err)
	}

	var quadrantConfig primitive.M

	err = json.Unmarshal([]byte(quadrantConfigResult), &quadrantConfig)

	if err != nil {
		fmt.Print("Error: ", err)
	}

	quadrantsConfigMap, err := quadrant.ConvertToQuadrantMap(quadrantConfig)
	if err != nil {
		log.Fatalf("Failed to convert quadrant config: %v", err)
	}

	quadrants := []*quadrant.Quadrant{
		quadrant.NewQuadrant(quadrantsConfigMap["QUADRANT_1"].(map[string]interface{})["Color"].(string), nil, "QUADRANT_1", quadrantsConfigMap["QUADRANT_1"].(map[string]interface{})["Path"].([]int)),
		quadrant.NewQuadrant(quadrantsConfigMap["QUADRANT_2"].(map[string]interface{})["Color"].(string), nil, "QUADRANT_2", quadrantsConfigMap["QUADRANT_2"].(map[string]interface{})["Path"].([]int)),
		quadrant.NewQuadrant(quadrantsConfigMap["QUADRANT_3"].(map[string]interface{})["Color"].(string), nil, "QUADRANT_3", quadrantsConfigMap["QUADRANT_3"].(map[string]interface{})["Path"].([]int)),
		quadrant.NewQuadrant(quadrantsConfigMap["QUADRANT_4"].(map[string]interface{})["Color"].(string), nil, "QUADRANT_4", quadrantsConfigMap["QUADRANT_4"].(map[string]interface{})["Path"].([]int)),
	}

	newBoard := CreateBoardInDB(boardId, autoPlay, playersRequiredToStartGame, ticketAmount, rakeAmount, rakeAmountType, autoPlayTimer)

	board := &Board{
		id:                         newBoard["boardId"].(string),
		quadrants:                  quadrants,
		players:                    []*player.Player{},
		safePositions:              quadrantsConfigMap["SafePositions"].([]int),
		ticketAmount:               ticketAmount,
		rakeAmount:                 rakeAmount,
		rakeAmountType:             rakeAmountType,
		autoPlay:                   autoPlay,
		autoPlayTimer:              autoPlayTimer,
		playersRequiredToStartGame: playersRequiredToStartGame,
		status:                     ludo_board_constants.BoardStatus("WAITING"),
	}

	return board
}

func (b *Board) GetMaxPlayers() int {
	return b.playersRequiredToStartGame
}

// GetID returns the unique identifier for the Board
func (b *Board) GetID() string {
	return b.id
}

func (b *Board) SetStatus(status ludo_board_constants.BoardStatus) {
	b.status = status
}

func (b *Board) GetSafePositions() []int {
	return b.safePositions
}

// GetBoardStatus returns the current game status
func (b *Board) GetBoardStatus() ludo_board_constants.BoardStatus {
	return b.status
}

func (b *Board) GetAutoPlay() bool {
	return b.autoPlay
}

func (b *Board) GetDiceRolledValue() int {
	return b.diceRolledValue
}

func (b *Board) SetDiceRolledValue(value int) {
	b.diceRolledValue = value
}

// GetPlayers returns all players in the game
func (b *Board) GetPlayers() []*player.Player {
	return b.players
}

func (b *Board) HasFinished() bool {
	return b.GetBoardStatus() == ludo_board_constants.FINISHED
}

// AddPlayer adds a new player to the game
func (b *Board) AddPlayer(playerId string, name string, walletAddress string) error {
	// log.Printf("Attempting to add player - ID: %s, Name: %s", playerId, name)

	if b.HasFinished() {
		return fmt.Errorf("game has already finished. Not accepting new players")
	}

	playerExists := b.GetPlayerByPlayerId(playerId)

	// log.Printf("Player exists: %+v", playerExists)

	//TODO: Create a handleReconnection function to handle reconnection
	if playerExists != nil {
		// log.Printf("[AddPlayer] Player with ID %s already exists in the game", playerId)
		err := b.HandleReconnection(playerExists)

		if err != nil {
			// log.Printf("[AddPlayer] Failed to handle reconnection for player %s", playerId)
			return err

		}

		// log.Printf("[AddPlayer] Successfully handled reconnection for player %s", playerId)
		return nil
	}

	// log.Printf("[AddPlayer] Checking board status: %s", b.GetBoardStatus())
	if b.GetBoardStatus() == ludo_board_constants.PLAYING {
		// log.Printf("[AddPlayer] Rejected: Game already started for board %s", b.GetID())
		return fmt.Errorf("game Already started. Not Accepting new players")
	}

	// log.Printf("[AddPlayer] Current player count: %d, Max players: %d", len(b.GetPlayers()), b.GetMaxPlayers())
	if len(b.GetPlayers()) == b.GetMaxPlayers() {
		// log.Printf("[AddPlayer] Rejected: Maximum players reached for board %s", b.GetID())
		return fmt.Errorf("game has reached maximum players. Not Accepting new players")
	}

	player := &player.Player{
		ID:               playerId,
		PlayerId:         playerId,
		Name:             name,
		ConnectionStatus: player.PLAYER_CONNECTED,
		WalletAddress:    walletAddress,
	}

	b.players = append(b.players, player)

	// log.Printf("[AddPlayer] Created new player instance - ID: %s, Name: %s", player.ID, player.Name)

	// log.Printf("[AddPlayer] Sending game initialization message to player %s", playerId)

	b.SendGameInitializeMessage(playerId)

	b.SendBoardWaitingPlayersMessage(player)

	if b.playersRequiredToStartGame == len(b.players) {
		log.Printf("[AddPlayer] All players have joined the board %s", b.GetID())
		log.Printf("[AddPlayer] Sending quadrant selection message")
		b.SendQuadrantSelectionMessage()
	}

	// log.Printf("[AddPlayer] Sending quadrant selection message")

	// log.Printf("[AddPlayer] Successfully added player %s to board %s. Total players: %d", playerId, b.GetID(), len(b.players))

	return nil
}

func (b *Board) HandleReconnection(existingPlayer *player.Player) error {
	if !existingPlayer.IsConnected() {
		// log.Printf("[AddPlayer] Player with ID %s reconnected to the game", playerId)
	}

	log.Printf(
		"[AddPlayer] Player %s reconnected to the game. Player details: ID: %s, Name: %s, Quadrant: %s",
	)
	existingPlayer.SetConnectionStatus(player.PLAYER_CONNECTED)
	log.Printf("[AddPlayer] Player %s reconnected to the game and connection status is %d", existingPlayer.GetPlayerId(), existingPlayer.ConnectionStatus)

	UpdatePlayerConnectionDetails(b.GetID(), existingPlayer.GetPlayerId(), "reconnection", time.Now())

	// log.Printf("[AddPlayer] Sending board reconnection message to player %s", playerId)
	boardReconnectionMessage := NewBoardReconnectionMessage(b.BuildBoardJoinedMessage().participants, b.GetPawnsPositionsInTheBoard())
	socket.SendMessage(existingPlayer.GetPlayerId(), boardReconnectionMessage, b.id)

	b.handleMessageAfterReconnection(*existingPlayer)

	return nil
}

func (b *Board) handleMessageAfterReconnection(existingPlayer player.Player) {
	if b.expectedMessage != nil && b.expectedMessage.PlayerId == existingPlayer.GetPlayerId() {
		if b.expectedMessage.EventName == ludo_board_constants.BOARD_MOVEPAWN {
			if quadrant := b.GetQuadrantFromPlayer(existingPlayer.GetPlayerId()); quadrant != nil && quadrant.GetName() == b.currentTurn {
				diceRolledMessage := dice.NewDiceRolledMessage(ludo_board_constants.BOARD_DICEROLLED, b.diceRolledValue, existingPlayer.Quadrant, b.calculateMovablePawns(b.diceRolledValue))
				b.SetExpectedMovePawnMessage(quadrant.GetName(), 30*time.Second, b.diceRolledValue)
				b.broadCastMessage(diceRolledMessage)
			}
		}

		if b.expectedMessage.EventName == ludo_board_constants.BOARD_DICEROLL || b.expectedMessage.EventName == ludo_board_constants.BOARD_TURN_COMPLETED {
			turnMessage := NewTurnMessage(ludo_board_constants.TURN, b.currentTurn, b.GetPawnsPositionsInTheBoard())
			b.SetExpectedDiceRollMessage(b.currentTurn, 30*time.Second)
			b.broadCastMessage(turnMessage)
		}

		if b.expectedMessage.EventName == ludo_board_constants.SELECT_QUADRANT {
			b.SendQuadrantSelectionMessage()
			b.SetExpectedQuadrantSelectMessage(b.expectedMessage.Quadrant, 30*time.Second)
		}
	}
}

func (b *Board) HandleDisconnection(playerId string) error {
	if b == nil {
		return errors.New("board instance is nil")
	}
	// log.Printf("[HandleDisconnection] Starting disconnection handling for player %s", playerId)
	// log.Printf("[HandleDisconnection] Current players in board %s: %+v", b.GetID(), b.GetPlayers())
	log.Println("Disconnected player is: ", b.GetPlayerByPlayerId(playerId))
	// Find and validate player's quadrant
	playerQuadrant := b.GetQuadrantFromPlayer(playerId)
	if playerQuadrant == nil {
		b.RemovePlayer(playerId)
		return fmt.Errorf("no quadrant found for player %s", playerId)
	}

	var disconnectedPlayer *player.Player
	allDisconnected := true
	var remainingPlayerId string
	remainingPlayersCount := 0

	for _, p := range b.GetPlayers() {
		// log.Printf("[HandleDisconnection] Checking player - Name: %s, ID: %s, Connected: %t", p.Name, p.ID, p.IsConnected())

		if p.ID == playerId {
			// log.Printf("[HandleDisconnection] Found disconnecting player - Name: %s, ID: %s", p.GetName(), p.GetPlayerId())
			disconnectedPlayer = p
			disconnectedPlayer.SetConnectionStatus(player.PLAYER_DISCONNECTED)
			UpdatePlayerConnectionDetails(b.GetID(), playerId, "disconnection", time.Now())
			// log.Printf("[HandleDisconnection] Marked player %s as disconnected", disconnectedPlayer.GetName())
		}

		if p.IsConnected() {
			allDisconnected = false
			remainingPlayerId = p.ID
			remainingPlayersCount++
			// log.Printf("[HandleDisconnection] Found connected player - ID: %s, Name: %s", p.GetPlayerId(), p.GetName())
		}
	}

	// log.Printf("[HandleDisconnection] Broadcasting disconnection message for player %s", playerId)
	disconnectionMessage := NewDisconnectionMessage(ludo_board_constants.PLAYER_DISCONNECTED, b.GetPlayerByPlayerId(playerId).Name)
	socket.BroadcastMessage(disconnectionMessage, b.GetID())

	if b.GetBoardStatus() != ludo_board_constants.FINISHED {
		// log.Printf("[HandleDisconnection] Game not finished, checking conditions")

		if !disconnectedPlayer.HasSelectedQuadrant() {
			// log.Printf("[HandleDisconnection] Player %s hadn't selected quadrant, removing from game", playerId)
			b.RemovePlayer(playerId)
			return nil
		}

		if disconnectedPlayer.HasSelectedQuadrant() && b.GetBoardStatus() == ludo_board_constants.WAITING {
			// log.Printf("[HandleDisconnection] Player %s had selected quadrant, game not started, removing from game", playerId)
			b.RemovePlayer(playerId)
			return nil
		}

		if allDisconnected {
			// log.Printf("[HandleDisconnection] All players disconnected, handling board cleanup")
			b.handleAllDisconnection()
		} else if remainingPlayersCount == 1 {
			// log.Printf("[HandleDisconnection] Only one player remaining (%s), starting 30-second timer", remainingPlayerId)
			go func() {
				time.Sleep(30 * time.Second)
				log.Printf("[HandleDisconnection] Timer completed, re-evaluating player count")

				connectedPlayersCount := 0
				for _, p := range b.GetPlayers() {
					log.Printf("[HandleDisconnection] Checking player %s - Connected: %t", p.GetName(), p.IsConnected())
					if p.IsConnected() {
						connectedPlayersCount++
					}
				}

				log.Printf("[HandleDisconnection] Connected players after timer: %d", connectedPlayersCount)
				if connectedPlayersCount == 1 {
					// log.Printf("[HandleDisconnection] Still only one player connected, ending game")
					b.handleAllDisconnectedExceptOne(remainingPlayerId)
				}
			}()
		} else {
			// log.Printf("[HandleDisconnection] Board %s continues with %d connected players", b.GetID(), remainingPlayersCount)
		}
	} else {
		// log.Printf("[HandleDisconnection] Game already finished, no further action needed")
	}

	return nil
}

func (b *Board) RemovePlayer(playerId string) {
	// Find player index and quadrant first
	var playerIndex int = -1
	var playerQuadrant *quadrant.Quadrant

	// Find player and their quadrant
	for i, p := range b.players {
		if p != nil && p.ID == playerId {
			playerIndex = i
			playerQuadrant = b.GetQuadrantFromPlayer(playerId)
			break
		}
	}

	// Only proceed if we found the player
	if playerIndex >= 0 {
		// Remove player from players slice
		b.players = append(b.players[:playerIndex], b.players[playerIndex+1:]...)

		// Clean up quadrant if it exists
		if playerQuadrant != nil {
			playerQuadrant.RemovePlayer() // This will handle setting isOccupied to false
		}
	}
}

func (b *Board) GetQuadrants() []*quadrant.Quadrant {
	return b.quadrants
}

func (b *Board) GetPlayerByPlayerId(playerId string) *player.Player {
	for _, player := range b.players {
		if player.GetPlayerId() == playerId {
			return player
		}
	}
	return nil
}

func (b *Board) GetPlayerByQuadrant(quadrant string) *player.Player {
	for _, player := range b.players {
		if player.GetQuadrant() == quadrant {
			return player
		}
	}
	return nil
}

func (q *Board) GetQuadrantFromPlayer(playerId string) *quadrant.Quadrant {
	for _, quadrant := range q.quadrants {
		if quadrant.GetPlayer() != nil && quadrant.GetPlayer().GetPlayerId() == playerId {
			return quadrant
		}
	}
	// log.Printf("No quadrant found for player %s", playerId)
	return nil
}

func (b *Board) IsBoardReadyToStartGame() bool {
	if len(b.players) < b.playersRequiredToStartGame {
		return false
	}
	for _, player := range b.players {
		if player.QuadrantSelectionStatus != 2 {
			return false
		}
	}
	return true
}

func (b *Board) SelectQuadrant(playerId string, selectQuadrantMessage quadrant.QuadrantSelectMessage) (*player.Player, error) {

	quadrantName := selectQuadrantMessage.GetQuadrant()

	for _, value := range b.GetAvailableQuadrants() {
		fmt.Printf("Key: %v, quadrant.Quadrant: %v\n", value, quadrantName)
		if value == quadrantName {

			err := CreateBetTransaction(b, playerId)

			if err != nil {
				newBoardBetFailedMessage := NewBoardBetFailedMessage(ludo_board_constants.BOARD_BET_FAILED, err.Error())
				socket.SendMessage(playerId, newBoardBetFailedMessage, b.id)
				b.RemovePlayer(playerId)
				return nil, err
			}
			// Create a new player with the selected color
			playerInstance := b.GetPlayerByPlayerId(playerId)

			quadrantInstance := b.GetQuadrant(quadrantName)

			quadrantInstance.Select(playerInstance)

			b.UpdateQuadrantSelection(playerId)

			newPlayer := &player.PlayerSchema{
				ID:       playerId,
				PlayerID: playerId,
				Quadrant: quadrantInstance.GetName(),
				Name:     playerInstance.Name,
				JoinedAt: time.Now(),
			}

			// Call the AddPlayerToDB function from the services package
			err = AddPlayerToBoardInDB(b.id, *newPlayer)

			if err != nil {
				return nil, err
			}

			newRoomJoinedMessage := b.BuildBoardJoinedMessage()

			b.broadCastMessage(newRoomJoinedMessage)

			b.startGameIfReady()

			if len(b.GetAvailableQuadrants()) > 0 {
				b.SendQuadrantSelectionMessage()

			}
			return playerInstance, nil
		}
	}
	return nil, errors.New("quadrant.Quadrant not available")
}

func (b *Board) startGameIfReady() {
	if b.IsBoardReadyToStartGame() {
		UpdateGameStatusAndAddStartTime(b.GetID(), ludo_board_constants.PLAYING, time.Now())

		b.SetFirstTurn()

		gameStartMessage := NewGameStartMessage(ludo_board_constants.GAME_START)

		b.SetStatus(ludo_board_constants.PLAYING)

		b.broadCastMessage(gameStartMessage)

		// gs := &LudoGameService{}
		// gs.CreateEmptyBoardInstances()

		turnMessage := NewTurnMessage(ludo_board_constants.TURN, b.GetFirstTurn(), b.GetPawnsPositionsInTheBoard())
		b.SetExpectedDiceRollMessage(b.GetFirstTurn(), 30*time.Second)
		b.broadCastMessage(turnMessage)
	}
}

func (b *Board) calculateMovablePawns(steps int) []string {
	var movablePawns []string
	quadrant := b.GetQuadrant(b.currentTurn)

	for _, p := range quadrant.GetPawns() {

		log.Printf("Checking if pawn %s is a valid move", p.GetName())

		if p.IsValidMove(steps) {
			log.Printf("pawn %s is a valid move with steps %d", p.GetName(), steps)
			movablePawns = append(movablePawns, p.GetName())

		}
	}

	return movablePawns
}

func (b *Board) BuildBoardJoinedMessage() *BoardJoinedMessage {
	var participants []ParticipantInfo

	// log the participants

	for _, player := range b.players {
		if player.QuadrantSelectionStatus != 2 {
			continue
		}
		participants = append(participants, ParticipantInfo{
			Player: Player{
				Id:   player.GetPlayerId(),
				Name: player.GetName(),
			},
			Quadrant: player.GetQuadrant(),
		})
	}

	// log.Println("Participants: ", participants)

	var playerSelectingTheQuadrant Player

	for _, player := range b.players {
		// log.Println("[BuildBoardJoinedMessage] Checking if player has turn to choose the quadrant", player.Name, player.QuadrantSelectionStatus)
		if player.QuadrantSelectionStatus == 1 {
			playerSelectingTheQuadrant = Player{
				Id:   player.GetPlayerId(),
				Name: player.GetName(),
			}
			break
		}
	}

	// log.Println("[BuildBoardJoinedMessage] Player selecting the quadrant: ", playerSelectingTheQuadrant)

	newRoomJoinedMessage := NewBoardJoinedMessage(
		ludo_board_constants.BOARD_JOINED,
		participants,
		playerSelectingTheQuadrant,
	)

	return newRoomJoinedMessage
}

func (b *Board) SetTicketAmount(ticketAmount int) {
	b.ticketAmount = ticketAmount
}

func (b *Board) GetTicketAmount() int {
	return b.ticketAmount
}

func (b *Board) SetRakeAmount(rakeAmount int) {
	b.rakeAmount = rakeAmount
}

func (b *Board) GetRakeAmount() int {
	return b.rakeAmount
}

func (b *Board) SetRakeAmountType(rakeAmountType ludo_board_constants.RakeAmountType) {
	b.rakeAmountType = rakeAmountType
}

func (b *Board) GetRakeAmountType() ludo_board_constants.RakeAmountType {
	return b.rakeAmountType
}

// GetAvailableColors returns the list of available colors
// Returns:
//   - map[int]string: Map of available colors
//
// returns only the string values
func (b *Board) GetAvailableQuadrants() []string {
	names := []string{}
	for _, value := range b.GetQuadrants() {
		if !value.GetIfQuadrantIsOccupied() {
			names = append(names, value.GetName())
		}
	}
	return names

}

// Set first turn
func (b *Board) SetFirstTurn() {
	for _, quadrant := range b.quadrants {
		if quadrant.GetPlayer() != nil {
			b.currentTurn = quadrant.GetName()
			b.nextTurn = b.quadrants[(1)%len(b.quadrants)].GetName()
			return
		}
	}
}

// GetFirstTurn returns the quadrant who will have the first turn will be the one with a player
// Returns:
//   - *quadrant.Quadrant: returns the quadrant name e.b: QUADRANT_1
func (b *Board) GetFirstTurn() string {
	for _, quadrant := range b.quadrants {
		if quadrant.GetPlayer() != nil {
			return quadrant.GetName()
		}
	}
	return "No player found"
}

// GetNextPlayer returns the player who will have the next turn
// Returns:
//   - *player.Player: Pointer to the next player
func (b *Board) GetNextTurn() string {
	return b.nextTurn
}

// GetQuadrant returns the quadrant corresponding to the specified color
// Parameters:
//   - color (string): The color of the quadrant ("red", "green", "yellow", "blue")
//
// Returns:
//   - *quadrant.Quadrant: Pointer to the quadrant.Quadrant if found, nil otherwise
func (b *Board) GetQuadrant(quadrantName string) *quadrant.Quadrant {
	for _, quadrant := range b.quadrants {
		if quadrant.GetName() == quadrantName {
			return quadrant
		}
	}
	return nil
}

// NextTurn manages the turn sequence between players based on the dice roll
// If a player rolls a 6, captures a pawn, or reaches the last position, they get another turn; otherwise, turn passes to the next player
// Parameters:
//   - diceRoll (int): The number rolled on the dice (1-6)
//   - captured (bool): Indicates if a pawn was captured
//   - reachedLastPosition (bool): Indicates if a pawn reached the last position
func (b *Board) NextTurn(diceRoll int, captured bool, reachedLastPosition bool, sendMessage bool) {
	if len(b.quadrants) == 0 {
		return
	}

	// If the dice roll is 6, a pawn was captured, or a pawn reached the last position, the current player retains their turn
	if diceRoll == 6 || captured || reachedLastPosition {
		// log.Printf("player.Player in quadrant %s gets another turn (diceRoll: %d, captured: %v, reachedLastPosition: %v)", b.currentTurn, diceRoll, captured, reachedLastPosition)
		if sendMessage {
			// Broadcast the next turn message
			turnMessage := NewTurnMessage(ludo_board_constants.TURN, b.currentTurn, b.GetPawnsPositionsInTheBoard())
			b.SetExpectedDiceRollMessage(b.currentTurn, 30*time.Second)
			b.broadCastMessage(turnMessage)
		}
		return
	}

	// Move to the next quadrant in the list
	currentIndex := -1
	for i, quadrant := range b.quadrants {
		if quadrant.GetName() == b.currentTurn {
			currentIndex = i
			break
		}
	}

	if currentIndex != -1 {
		for i := 1; i < len(b.quadrants); i++ {
			nextIndex := (currentIndex + i) % len(b.quadrants)
			if b.quadrants[nextIndex].GetPlayer() != nil {
				b.currentTurn = b.quadrants[nextIndex].GetName()
				b.nextTurn = b.quadrants[(nextIndex+1)%len(b.quadrants)].GetName()
				log.Printf("Next turn is for quadrant %s", b.currentTurn)
				break
			}
		}
	}

	if sendMessage {
		// Broadcast the next turn message
		turnMessage := NewTurnMessage(ludo_board_constants.TURN, b.currentTurn, b.GetPawnsPositionsInTheBoard())
		b.SetExpectedDiceRollMessage(b.currentTurn, 30*time.Second)
		b.broadCastMessage(turnMessage)
	}
}

// TurnCompleted broadcasts a message to indicate that the current player's turn has been completed
func (b *Board) TurnCompleted() {
	turnMessage := NewTurnMessage(ludo_board_constants.TURN, b.currentTurn, b.GetPawnsPositionsInTheBoard())
	b.SetExpectedDiceRollMessage(b.currentTurn, 30*time.Second)
	b.broadCastMessage(turnMessage)
}

// Get all the pawns' positions in the board
func (b *Board) GetPawnsPositionsInTheBoard() []pawn.PawnPositions {
	var allPawnPositions []pawn.PawnPositions

	for _, quadrant := range b.quadrants {
		pawnPositions := pawn.PawnPositions{
			Quadrant: quadrant.GetName(),
		}

		for _, p := range quadrant.GetPawns() {
			position := p.GetPosition()
			positionValue := -1
			if position != nil {
				positionValue = *position
			}
			// log.Printf("Pawn %s in quadrant %s is at position %v", p.GetName(), quadrant.GetName(), positionValue)
			pawnPositions.AddPosition(p.GetName(), positionValue)
		}

		allPawnPositions = append(allPawnPositions, pawnPositions)
	}

	// log.Printf("All pawn positions: %v", allPawnPositions)
	return allPawnPositions
}

// MovePawn executes a pawn movement for a player by the specified number of steps
// Parameters:
//   - pawnMoveMessage (messages.PawnMoveMessage): The message containing the details of the pawn move
//
// Returns:
//   - error: Error if the move is invalid, nil if successful
func (b *Board) MovePawn(pawnMoveMessage pawn.PawnMoveMessage) error {
	quadrantName := pawnMoveMessage.GetQuadrant()
	pawnName := pawnMoveMessage.GetPawn()
	steps := pawnMoveMessage.GetSteps()
	if b.expectedMessage != nil {
		if b.expectedMessage.Steps != steps {
			err := fmt.Errorf("invalid steps %d, expected %d on quadrant %s by player %s", steps, b.expectedMessage.Steps, quadrantName, b.GetPlayerByQuadrant(quadrantName).GetPlayerId())
			// log.Println(err)
			return err
		}
	}

	// log.Printf("Attempting to move pawn %s in quadrant %s by %d steps", pawnName, quadrantName, steps)

	// Get the player's quadrant
	quadrantInstance := b.GetQuadrant(quadrantName)

	if quadrantInstance == nil {
		err := fmt.Errorf("no quadrant found with name %s", quadrantName)
		// log.Println(err)
		return err
	}

	// Get the player from the quadrant
	player := quadrantInstance.GetPlayer()

	if player == nil {
		err := fmt.Errorf("no player found in quadrant %s", quadrantName)
		// log.Println(err)
		return err
	}

	// Get the pawn by name
	pawnInstance := quadrantInstance.GetPawnByName(pawnName)

	if pawnInstance == nil {
		err := fmt.Errorf("pawn %s not found for player %s", pawnName, player.ID)
		// log.Println(err)
		return err
	}

	// Move the pawn
	moveResult, err := pawnInstance.MovePawn(steps, quadrantName)

	if err != nil {
		// log.Printf("Failed to move pawn %s: %v", pawnName, err)
		return err
	}

	// Handle capturing of opponent's pawn
	capturedPawns := b.capturePawnIfPresent(pawnInstance, quadrantInstance)

	// Update the pawn movement in the database
	boardDAO := NewBoardDAO()

	err = boardDAO.UpdateBoardPawnMovement(b.id, quadrantName, pawnName, moveResult.InitialPosition, moveResult.FinalPosition, time.Now(), steps)

	if err != nil {
		// log.Printf("Failed to update pawn movement in database: %v", err)
		return err
	}

	// Check if all pawns of the player are finished
	allFinished := true

	for _, p := range quadrantInstance.GetPawns() {
		if !p.IsAtFinish() {
			allFinished = false
			break
		}
	}

	movementDetails := pawn.NewPawnMovedMessage(map[string]interface{}{
		"eventName": ludo_board_constants.BOARD_PAWNMOVED,
		"pawn":      pawnInstance.GetName(),
		"steps":     steps,
		"responseCode": func() int {
			if moveResult.ValidationErrors[0].Message != "" {
				return 400
			} else {
				return 200
			}
		}(),
		"quadrant":         quadrantInstance.GetName(),
		"initialPosition":  moveResult.InitialPosition,
		"finalPosition":    moveResult.FinalPosition,
		"initialIndex":     moveResult.InitialIndex,
		"finalIndex":       moveResult.FinalIndex,
		"isAtHome":         moveResult.IsAtHome,
		"capturedPawns":    capturedPawns,
		"validationErrors": moveResult.ValidationErrors,
		"positions":        b.GetPawnsPositionsInTheBoard(),
	})

	// Broadcast the movement details
	b.broadCastMessage(movementDetails)

	// Emit Game.Winner message if all pawns are finished
	if allFinished {
		// log.Printf("Player %s has finished all pawns. Emitting Game.Winner message.", player.GetPlayerId())
		// Update game status to completed
		err = UpdateBoardStatusAndAddEndTimeInDB(b.GetID(), ludo_board_constants.FINISHED, time.Now())
		if err != nil {
			// log.Printf("Failed to update game status and end time in database: %v", err)
		}

		// Update game winner
		err = UpdateGameWinnerInDB(b.GetID(), player.GetPlayerId(), b.getWinningAmount())

		if err != nil {
			log.Printf("Failed to update game winner in database: %v", err)
		}

		b.SetStatus(ludo_board_constants.FINISHED)

		endMessage := NewGameEndMessage(ludo_board_constants.GAME_END, player.GetPlayerId(), b.getWinningAmount(), 200)

		b.broadCastMessage(endMessage)

		CreateWinTransaction(b, player.GetPlayerId(), b.getWinningAmount())

		// gameEndMessage := NewGameEndMessage(ludo_board_constants.GAME_END)

		b.broadCastMessage(endMessage)

		return nil
	}

	// log.Printf("Broadcasted movement details for pawn %s", pawnName)

	// Handle next turn
	b.NextTurn(steps, len(capturedPawns) > 0, moveResult.IsAtHome, false)

	b.SetExpectedTurnCompletedMessage(quadrantInstance.GetName(), 30*time.Second)

	// log.Printf("Broadcasted turn message for quadrant %s", b.currentTurn)

	return nil
}

// IsGameOver checks if any player has won the game
// Returns:
//   - bool: True if the game is over, false otherwise
func (b *Board) IsGameOver() bool {
	for _, quadrant := range b.quadrants {
		if quadrant.HasWon() {
			return true
		}
	}
	return false
}

// capturePawnIfPresent checks if any opponent's pawn is at the same position after the move and captures it
// Parameters:
//   - pawn (*pawn.Pawn): The pawn that just moved
//   - player (*player.Player): The player who moved the pawn
//
// Returns:
//   - []string: List of captured pawn names
func (b *Board) capturePawnIfPresent(pawn *pawn.Pawn, quadrant *quadrant.Quadrant) []string {
	capturedPawns := []string{}
	// log.Printf("Checking for captures by pawn %s of quadrant %s at position %v", pawn.GetName(), quadrant.GetName(), pawn.GetPosition())
	for _, otherQuadrant := range b.quadrants {
		if otherQuadrant.GetName() == quadrant.GetName() {
			continue
		}

		if otherQuadrant == nil {
			// log.Printf("Opponent quadrant %s is nil", otherQuadrant.GetName())
			continue
		}
		for _, otherQuadrantPawn := range otherQuadrant.GetPawns() {

			if otherQuadrantPawn == nil {
				// log.Printf("Opponent pawn is nil for quadrant %s", otherQuadrant.GetName())
				continue
			}

			// log.Printf("Checking opponent pawn %s of quadrant %s at position %v", otherQuadrantPawn.GetName(), otherQuadrant.GetName(), otherQuadrantPawn.GetPosition())

			if otherQuadrantPawn.GetPosition() != nil && pawn.GetPosition() != nil && *otherQuadrantPawn.GetPosition() == *pawn.GetPosition() && !otherQuadrantPawn.IsIdle() {
				// Check if the position is a safe zone
				isSafeZone := false
				for _, safePos := range b.safePositions {
					if safePos == *pawn.GetPosition() {
						isSafeZone = true
						break
					}
				}
				// Capture the opponent's pawn if not in a safe zone
				if !isSafeZone {
					// log.Printf("Capturing pawn %s of player %s at position %d", otherQuadrantPawn.GetName(), otherQuadrant.GetName(), *otherQuadrantPawn.GetPosition())
					otherQuadrantPawn.SetPosition(nil)
					otherQuadrantPawn.SetStatus(ludo_board_constants.PAWN_IDLE)
					otherQuadrantPawn.SetPosition(nil) // Reset the position to nil
					capturedPawns = append(capturedPawns, otherQuadrantPawn.GetName())
					// log.Printf("player.Player %s's pawn %s captured player %s's pawn %s!", quadrant.GetName(), pawn.GetName(), otherQuadrant.GetName(), otherQuadrantPawn.GetName())
				} else {
					// log.Printf("pawn.Pawn %s of player %s is in a safe zone at position %d and cannot be captured", otherQuadrantPawn.GetName(), otherQuadrant.GetName(), *otherQuadrantPawn.GetPosition())
				}
			} else {
				// log.Printf("No capture: Opponent pawn %s of player %s is at position %v, current pawn %s of player %s is at position %v",
				// 	otherQuadrantPawn.GetName(), otherQuadrant.GetName(), otherQuadrantPawn.GetPosition(), pawn.GetName(), quadrant.GetName(), pawn.GetPosition())
			}
		}
	}
	return capturedPawns
}

// DiceRoll handles the dice roll for a player
// Parameters:
//   - playerId (string): The ID of the player rolling the dice
func (b *Board) DiceRoll(playerId string) {
	// log.Printf("player.Player %s rolled the dice", playerId)

	diceRollingMessage := dice.NewDiceRollingMessage(ludo_board_constants.BOARD_DICEROLLING)
	b.broadCastMessage(diceRollingMessage)

	quadrantInstance := b.GetQuadrantFromPlayer(playerId)

	if quadrantInstance == nil {
		// log.Printf("No quadrant found for player %s", playerId)
		return
	}

	// Check if the player has any pawns unlocked
	hasUnlockedPawns := false

	for _, pawn := range quadrantInstance.GetPawns() {
		if !pawn.IsIdle() {
			hasUnlockedPawns = true
			break
		}
	}

	diceInstance := &dice.Dice{}
	diceValue := diceInstance.Roll()
	// log.Printf("Dice value: %d", diceValue)
	time.Sleep(300 * time.Millisecond)

	b.SetDiceRolledValue(diceValue)

	movablePawns := b.calculateMovablePawns(diceValue)

	// Broadcast the dice rolled message once
	diceRolledMessage := dice.NewDiceRolledMessage("Board.DiceRolled", diceValue, quadrantInstance.GetName(), movablePawns)
	b.SetExpectedMovePawnMessage(quadrantInstance.GetName(), 30*time.Second, diceValue)
	b.broadCastMessage(diceRolledMessage)
	// log.Printf("Broadcasted dice rolled message")

	// If the player has no pawns unlocked and the dice value is not 6, pass the turn
	if !hasUnlockedPawns && diceValue != 6 {
		// log.Printf("player.Player %s has no unlocked pawns and rolled a %d. Passing turn to the next player.", playerId, diceValue)
		b.NextTurn(diceValue, false, false, true)
		return
	}

	// Check if any pawn can move the rolled dice value
	canMove := false
	for _, pawn := range quadrantInstance.GetPawns() {
		if (pawn.IsIdle() && diceValue == 6) || (!pawn.IsIdle() && pawn.IsValidMove(diceValue)) {
			canMove = true
			break
		}
	}

	// If no pawn can move the rolled dice value, pass the turn
	if !canMove {
		// log.Printf("player.Player %s has no pawns that can move %d steps. Passing turn to the next player.", playerId, diceValue)
		b.NextTurn(diceValue, false, false, true)
		return
	}

	// If the dice roll is 6, the player gets another turn
	if diceValue == 6 {
		// log.Printf("player.Player %s rolled a 6 and gets another turn.", playerId)
		return
	}

	// If the dice roll is not 6 and the player has pawns that can move, no turn message is broadcasted
	// log.Printf("player.Player %s has pawns that can move %d steps.", playerId, diceValue)
}

func (b *Board) UpdateQuadrantSelection(playerId string) {
	// log.Printf("Updating quadrant selection for playerId %s", playerId)
	for _, player := range b.players {
		// log.Printf("Player %s quadrant selection status: %d", player.ID, player.QuadrantSelectionStatus)
		if player.QuadrantSelectionStatus == 1 && player.ID == playerId {
			player.QuadrantSelectionStatus = 2
			// log.Printf("Updated player.Player %s quadrant selection", player.ID)
			// log.Printf("Quadrant selection status: %d is for Player: %s", player.QuadrantSelectionStatus, player.PlayerId)
			return
		}
	}
}

func (b *Board) SendQuadrantSelectionMessage() {
	var playerToSend *player.Player

	for _, player := range b.players {
		// log.Printf("Player %s's Quadrant selection status: %d", player.ID, player.QuadrantSelectionStatus)
		if player.QuadrantSelectionStatus == 1 {
			// log.Printf("player.Player %s already in process of selection", player.GetPlayerId())
			return
		}
		if player.QuadrantSelectionStatus == 0 {
			playerToSend = player
			break
		}
	}

	if playerToSend == nil {
		// log.Printf("All players have selected their quadrants")
		return
	}

	// log.Printf("Sending quadrant selection message to player %s", playerToSend.ID)

	playerToSend.QuadrantSelectionStatus = 1

	quadrantSelectionPromptMessage := quadrant.NewSelectQuadrantMessage(ludo_board_constants.SELECT_QUADRANT, b.GetAvailableQuadrants(), 200)
	b.SetExpectedQuadrantSelectMessage(playerToSend.GetQuadrant(), 30*time.Second)

	boardSelectingQuadrantMessage := NewBoardSelectingQuadrantMessage(
		ludo_board_constants.BOARD_SELECTING_QUADRANT,
		Player{
			Id:   playerToSend.GetPlayerId(),
			Name: playerToSend.GetName(),
		},
	)

	b.broadCastMessage(boardSelectingQuadrantMessage)
	socket.SendMessage(playerToSend.ID, quadrantSelectionPromptMessage, b.id)
}

func (b *Board) SendGameInitializeMessage(playerId string) {

	safePositions := b.GetSafePositions()

	var quadrants []Quadrant

	for _, quadrant := range b.GetQuadrants() {
		quadrants = append(quadrants, Quadrant{
			Name:  quadrant.GetName(),
			Color: quadrant.GetColor(),
			Pawns: quadrant.GetPawnNames(),
			Path:  quadrant.GetPath(),
		})
	}

	var playerSelectingTheQuadrant Player

	for _, player := range b.GetPlayers() {
		if player.QuadrantSelectionStatus == 1 {
			playerSelectingTheQuadrant = Player{
				Id:   player.GetPlayerId(),
				Name: player.GetName(),
			}
			break
		}
	}

	gameInitializeMessage := NewGameInitializeMessage(ludo_board_constants.GAME_INITIALIZE, safePositions, quadrants, b.autoPlay, b.playersRequiredToStartGame, b.ticketAmount, playerSelectingTheQuadrant, b.autoPlayTimer)

	socket.SendMessage(playerId, gameInitializeMessage, b.id)

}

func (b *Board) SendBoardWaitingPlayersMessage(player *player.Player) {
	var playerSelectingQuadrant Player
	var waitingPlayers []Player

	newPlayer := Player{
		Id:   player.GetPlayerId(),
		Name: player.GetName(),
	}

	log.Printf("New player: %v", newPlayer)

	for _, player := range b.GetPlayers() {

		waitingPlayers = append(waitingPlayers, Player{
			Id:   player.GetPlayerId(),
			Name: player.GetName(),
		})

		if player.QuadrantSelectionStatus == 1 {
			playerSelectingQuadrant = Player{
				Id:   player.GetPlayerId(),
				Name: player.GetName(),
			}
		}
	}

	log.Printf("Waiting players: %v", waitingPlayers)
	log.Printf("Player selecting quadrant: %v", playerSelectingQuadrant)

	boardWaitingPlayersMessage := NewBoardWaitingPlayersMessage(ludo_board_constants.BOARD_WAITING_PLAYERS, waitingPlayers,
		newPlayer, playerSelectingQuadrant)

	socket.BroadcastMessage(boardWaitingPlayersMessage, b.GetID())
}

func (b *Board) getWinningAmount() int {
	// log.Printf("Calculating winning amount with ticket amount: %d, players required: %d, rake amount: %d, rake amount type: %s",
	// b.ticketAmount, b.playersRequiredToStartGame, b.rakeAmount, b.rakeAmountType)

	poolAmount := b.ticketAmount * b.playersRequiredToStartGame

	// winningAmount := b.ticketAmount * b.playersRequiredToStartGame
	// log.Printf("Total prize pool before rake: %d", poolAmount)

	var winningAmount int

	if b.rakeAmountType == ludo_board_constants.PERCENTAGE && b.rakeAmount > 0 {
		winningAmount = poolAmount - (poolAmount * b.rakeAmount / 100)
	} else {
		// log.Printf("Using fixed rake amount: %d", b.rakeAmount)
		winningAmount = poolAmount - b.rakeAmount
	}

	// log.Printf("Final winning amount after deducting rake: %d", winningAmount)

	return winningAmount
}

func (b *Board) handleAllDisconnection() error {
	// log.Printf("All players disconnected, discarding board %s", b.GetID())

	for _, p := range b.GetPlayers() {
		b.CreateRefundTransaction(p.GetPlayerId(), float64(b.GetTicketAmount()))
		b.GetQuadrant(p.GetQuadrant()).RemovePlayer()
		b.RemovePlayer(p.ID)
	}

	b.SetStatus(ludo_board_constants.DISCARDED)

	err := UpdateBoardStatusAndAddEndTimeInDB(b.GetID(), ludo_board_constants.DISCARDED, time.Now())

	if err != nil {
		// log.Printf("Failed to update game status and end time in database: %v", err)
	}

	return nil
}

func (b *Board) handleAllDisconnectedExceptOne(remainingPlayerId string) error {

	remainingPlayer := b.GetPlayerByPlayerId(remainingPlayerId)

	err := UpdateBoardStatusAndAddEndTimeInDB(b.GetID(), ludo_board_constants.FINISHED, time.Now())
	if err != nil {
		// log.Printf("Failed to update game status and end time in database: %v", err)
	}

	// Update game winner
	err = UpdateGameWinnerInDB(b.GetID(), remainingPlayer.GetPlayerId(), b.getWinningAmount())

	if err != nil {
		// log.Printf("Failed to update game winner in database: %v", err)
	}

	b.SetStatus(ludo_board_constants.FINISHED)

	CreateWinTransaction(b, remainingPlayerId, b.getWinningAmount())

	endMessage := NewGameEndMessage(ludo_board_constants.GAME_END, remainingPlayer.GetPlayerId(), b.getWinningAmount(), 200)

	socket.BroadcastMessage(endMessage, b.GetID())

	for _, p := range b.GetPlayers() {
		b.RemovePlayer(p.ID)
		// log.Printf("Removed player %s from board %s as the winner has been declared", p.ID, b.GetID())
	}
	// delete(Bs, b.GetID())
	// log.Printf("Board %s ended as only one player is left", b.GetID())

	return nil
}

func UpdateGameStatusAndAddStartTime(boardId string, status ludo_board_constants.BoardStatus, sTime time.Time) error {

	BoardDAO := NewBoardDAO()

	_, error := primitive.ObjectIDFromHex(boardId)

	if error != nil {
		log.Fatal(error)
	}

	err := BoardDAO.UpdateGameStatusAndAddStartTime(boardId, status, sTime)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func CreateBoardInDB(boardId string, autoPlay bool, playersRequiredToStartGame int, ticketAmount int, rakeAmount int, rakeAmountType ludo_board_constants.RakeAmountType, autoPlayTimer int) primitive.M {

	BoardDAO := NewBoardDAO()

	game := BoardSchema{
		ID:                         primitive.NewObjectID().Hex(),
		BoardId:                    boardId,
		Status:                     ludo_board_constants.WAITING,
		AutoPlay:                   autoPlay,
		AutoPlayTimer:              autoPlayTimer,
		TicketAmount:               ticketAmount,
		RakeAmount:                 rakeAmount,
		RakeAmountType:             rakeAmountType,
		PlayersRequiredToStartGame: playersRequiredToStartGame,
		StartTime:                  nil,
		EndTime:                    nil,
		Winner:                     nil,
		Players:                    []player.PlayerSchema{},
		PawnMoves:                  make(map[string]map[string][]MoveSchema),
	}

	result, err := BoardDAO.InsertBoard(game)

	if err != nil {
		log.Fatal(err)
	}

	return result
}

func AddPlayerToBoardInDB(boardId string, player player.PlayerSchema) error {
	// Retrieve the game by boardId
	game, err := NewBoardDAO().GetBoardById(boardId)

	if err != nil {
		log.Fatalf("failed to get game with ID %s: %v", boardId, err)
		return err
	}

	err = NewBoardDAO().AddPlayerToBoard(game.BoardId, player)

	if err != nil {
		return fmt.Errorf("failed to save game with ID %s: %v", game.ID, err)
	}
	return nil
}

func UpdatePlayerConnectionDetails(boardId, playerId string, cType string, time time.Time) error {
	boardDAO := NewBoardDAO()

	err := boardDAO.UpdatePlayerConnectionDetails(boardId, playerId, cType, time)

	if err != nil {
		log.Printf("Failed to update player disconnection time in database: %v", err)

		return err
	}

	return nil
}

func UpdateBoardStatusAndAddEndTimeInDB(boardId string, status ludo_board_constants.BoardStatus, endTime time.Time) error {
	boardDAO := NewBoardDAO()
	err := boardDAO.UpdateBoardStatusAndAddEndTime(boardId, status, endTime)
	if err != nil {
		log.Printf("Failed to update game status and end time in database: %v", err)
		return err
	}
	return nil
}

func UpdateGameWinnerInDB(boardId, winner string, winningAmount int) error {
	boardDAO := NewBoardDAO()

	err := boardDAO.UpdateGameWinner(boardId, winner, winningAmount)
	if err != nil {
		log.Printf("Failed to update game winner in database: %v", err)
		return err
	}

	return nil
}

func CreateBetTransaction(board *Board, playerId string) error {

	player := board.GetPlayerByPlayerId(playerId)

	payload := BetAndWinPayload{
		WalletAddress: player.WalletAddress,
		Amount:        float64(board.ticketAmount),
	}

	// player.SetBetId(payload.BetId)

	var response BetResponse
	err := createTransaction(board, playerId, "/core/crypto/game/ludoBet", payload, &response)
	if err != nil {
		return err
	}

	if response.Code != "RS200" {
		// log.Printf("Bet transaction failed with code: %s", response.Code)
		switch response.Code {
		case "RS405":
			return fmt.Errorf(ludo_board_constants.RS405)
		}
	}

	// log.Printf("Bet transaction completed successfully for player %s", playerId)
	return nil
}

func CreateWinTransaction(board *Board, playerId string, winningAmount int) error {

	player := board.GetPlayerByPlayerId(playerId)

	payload := BetAndWinPayload{
		WalletAddress: player.WalletAddress,
		Amount:        float64(winningAmount),
	}

	var response WinResponse

	err := createTransaction(board, playerId, "/core/crypto/game/ludoWin", payload, &response)
	if err != nil {
		return err
	}

	if response.Code != "RS200" {
		// log.Printf("Win transaction failed with code: %s", response.Code)
		return fmt.Errorf("win transaction failed with code: %s", response.Code)
	}

	// log.Printf("Win transaction completed successfully for player %s", playerId)
	return nil
}

func (b *Board) CreateRefundTransaction(playerId string, amount float64) error {
	payload := RefundPayload{
		PlayerId:        playerId,
		Amount:          amount,
		TransactionUuid: uuid.New().String(),
		RequestUuid:     uuid.New().String(),
		Currency:        "INR",
		GameId:          b.GetID(),
	}

	var response RefundResponse
	err := createTransaction(b, playerId, "/wallet/refund", payload, &response)
	if err != nil {
		return err
	}

	// log.Printf("Refund transaction completed successfully for player %s", playerId)
	return nil
}

func (b *Board) broadCastMessage(msg common.Message) {
	socket.BroadcastMessage(msg, b.GetID())
}

func createTransaction(board *Board, playerId string, endpoint string, payload interface{}, response interface{}) error {
	cfg := config.GetConfig()
	// log.Printf("Creating transaction for player %s with ticket amount %d", playerId, board.ticketAmount)

	if board.ticketAmount == 0 {
		// log.Printf("Skipping transaction: ticket amount is 0")
		return nil
	}

	// log.Printf("Created payload: %+v", payload)

	jsonData, err := json.Marshal(payload)
	if err != nil {
		// log.Printf("Failed to marshal payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %v", err)
	}
	// log.Printf("Marshalled JSON data: %s", string(jsonData))

	// log.Printf("Making POST request to %s", cfg.BasePlatformAPIUrl+endpoint)
	resp, err := http.Post(cfg.BasePlatformAPIUrl+endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		// log.Printf("Failed to make transaction request: %v", err)
		return err
	}

	defer resp.Body.Close()
	// log.Printf("Received response with status code: %d", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// log.Printf("Failed to read response body: %v", err)
		return fmt.Errorf("failed to read response body: %v", err)
	}
	// log.Printf("Response body: %s", string(body))

	if err := json.Unmarshal(body, &response); err != nil {
		// log.Printf("Failed to unmarshal response: %v", err)
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}
	// log.Printf("Parsed response: %+v", response)

	if resp.StatusCode != http.StatusOK {
		// log.Printf("Transaction failed with HTTP status: %d", resp.StatusCode)
		return fmt.Errorf("transaction failed with status: %d", resp.StatusCode)
	}

	return nil
}

type BetAndWinPayload struct {
	WalletAddress string  `json:"walletAddress"`
	Amount        float64 `json:"amount"`
}

type RefundPayload struct {
	PlayerId        string  `json:"playerId"`
	Amount          float64 `json:"amount"`
	TransactionUuid string  `json:"transactionUuid"`
	RequestUuid     string  `json:"requestUuid"`
	Currency        string  `json:"currency"`
	GameId          string  `json:"gameId"`
}

type BetResponse struct {
	User        string  `json:"user"`
	Code        string  `json:"code"`
	Status      string  `json:"status"`
	RequestUuid string  `json:"requestUuid"`
	Balance     float64 `json:"balance"`
}

type RefundResponse struct {
	User        string  `json:"user"`
	Status      string  `json:"status"`
	RequestUuid string  `json:"requestUuid"`
	Balance     float64 `json:"balance"`
}

type WinResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		PlayerResponse struct {
			ID          string   `json:"_id"`
			PhoneNumber string   `json:"phoneNumber"`
			SkinID      string   `json:"skinId"`
			Tags        []string `json:"tags"`
			Balance     []struct {
				// Define fields for balance object
			} `json:"balance"`
			BankDetails []struct {
				// Define fields for bank details array
			} `json:"bankDetails"`
			Email          string  `json:"email"`
			Name           string  `json:"name"`
			CurrentBalance float64 `json:"currentBalance"`
		} `json:"playerResponse"`
		NewTransaction struct {
			Player          string  `json:"player"`
			Amount          float64 `json:"amount"`
			SkinID          string  `json:"skinId"`
			Currency        string  `json:"currency"`
			TransactionType string  `json:"transactionType"`
			MoneyType       string  `json:"moneyType"`
			OpeningBalance  float64 `json:"openingBalance"`
			ClosingBalance  float64 `json:"closingBalance"`
			Details         []struct {
				// Define fields for details object
			} `json:"details"`
			IsTransactionSuccess bool   `json:"isTransactionSuccess"`
			ID                   string `json:"_id"`
			CreatedAt            string `json:"createdAt"`
			UpdatedAt            string `json:"updatedAt"`
			Version              int    `json:"__v"`
		} `json:"newTransaction"`
	} `json:"data"`
}

func (b *Board) GetExpectedMessage() *ExpectedMessage {
	return b.expectedMessage
}

func (b *Board) UnsetExpectedMessage() {
	b.expectedMessage = nil
}

func (b *Board) SetExpectedDiceRollMessage(quadrant string, timeout time.Duration) {
	// log.Printf("Setting expected dice roll message for quadrant %s", quadrant)
	b.expectedMessage = &ExpectedMessage{
		EventName: ludo_board_constants.BOARD_DICEROLL,
		Quadrant:  quadrant,
		PlayerId:  b.GetPlayerByQuadrant(quadrant).GetPlayerId(),
		TStamp:    time.Now(),
		Timeout:   timeout,
	}
}

func (b *Board) SetExpectedQuadrantSelectMessage(quadrant string, timeout time.Duration) {
	// log.Printf("Setting expected quadrant select message for quadrant %s", quadrant)
	b.expectedMessage = &ExpectedMessage{
		EventName: ludo_board_constants.QUADRANT_SELECT,
		Quadrant:  quadrant,
		PlayerId:  b.GetPlayerByQuadrant(quadrant).GetPlayerId(),
		TStamp:    time.Now(),
		Timeout:   timeout,
	}
}

func (b *Board) SetExpectedMovePawnMessage(quadrant string, timeout time.Duration, steps int) {
	// log.Printf("Setting expected move pawn message for quadrant %s", quadrant)
	b.expectedMessage = &ExpectedMessage{
		EventName: ludo_board_constants.BOARD_MOVEPAWN,
		Quadrant:  quadrant,
		PlayerId:  b.GetPlayerByQuadrant(quadrant).GetPlayerId(),
		TStamp:    time.Now(),
		Timeout:   timeout,
		Steps:     steps,
	}
}

func (b *Board) SetExpectedTurnCompletedMessage(quadrant string, timeout time.Duration) {
	// log.Printf("[@SetExpectedTurnCompletedMessage] Setting expected turn completed message for quadrant %s", quadrant)
	b.expectedMessage = &ExpectedMessage{
		EventName: ludo_board_constants.BOARD_TURN_COMPLETED,
		Quadrant:  quadrant,
		PlayerId:  b.GetPlayerByQuadrant(quadrant).GetPlayerId(),
		TStamp:    time.Now(),
		Timeout:   timeout,
	}
}
