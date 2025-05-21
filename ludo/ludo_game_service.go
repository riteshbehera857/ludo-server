package ludo

import (
	"fmt"
	"log"
	"ludo/board"
	"ludo/ludo_board_constants"
	"ludo/quadrant"
	"messaging/common"
	"time"

	"ludo/pawn"

	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var BoardInstances map[string]*board.Board = make(map[string]*board.Board)

type BoardConfig struct {
	boardId        string
	playerCount    int
	rakeAmountType ludo_board_constants.RakeAmountType
	amount         int
}

type LudoGameService struct {
	common.GameService
}

func (gs *LudoGameService) ProcessMessage(boardId string, playerId string, message common.Message, rawBytes []byte) error {
	var socketMessage = message.(common.SocketMessage)

	// log.Printf("Player %s sent message %s", playerId, socketMessage.GetEventName())

	eventParts := strings.Split(socketMessage.GetEventName(), ".")

	if len(eventParts) != 2 {
		log.Printf("Invalid message %s from player %s", socketMessage.GetEventName(), playerId)
		return fmt.Errorf("invalid event name format: %s", socketMessage.GetEventName())
	}

	className := eventParts[0]
	methodName := eventParts[1]

	board, exists := BoardInstances[boardId]

	if !exists {
		// log.Printf("Game instance not found for room ID: %s", boardId)
		return fmt.Errorf("game instance not found for room ID: %s", boardId)
	}

	playerObj := board.GetPlayerByPlayerId(playerId)

	expectedMessage := board.GetExpectedMessage()

	// Log the expected message
	// log.Println("[@ProcessMessage] Expected message: ", expectedMessage)

	if expectedMessage != nil {
		obj := *expectedMessage
		if obj.EventName != socketMessage.GetEventName() {
			err := fmt.Errorf("invalid event name : %s", socketMessage.GetEventName())
			log.Printf("Error: %s", err)
			return err
		}

		if obj.PlayerId != playerId {
			err := fmt.Errorf("Invalid player %s, expected message from player %s", playerId, obj.PlayerId)
			log.Printf("Error: %s", err)
			return err
		}

		if obj.Quadrant != playerObj.GetQuadrant() {
			err := fmt.Errorf("Invalid quadrant %s from player %s", playerObj.GetQuadrant(), obj.Quadrant)
			log.Printf("Error: %s", err)
			return err
		}
		// Unset the expected message once it is received
		board.UnsetExpectedMessage()
	} else {
		log.Printf("No expected message for event: %s, player: %s, quadrant: %s", socketMessage.GetEventName(), playerId, playerObj.GetQuadrant())
	}

	// log.Printf("Invoking method %s on class %s", methodName, className)

	var instance = board

	var args []reflect.Value

	switch methodName {
	case string(ludo_board_constants.SELECT):
		selectQuadrantMessage := &quadrant.QuadrantSelectMessage{}
		selectQuadrantMessageObject, err := selectQuadrantMessage.ToObject(string(rawBytes))

		// fmt.Println(selectQuadrantMessageObject)

		if err != nil {
			log.Printf("Error %s when parsing pawn move message", err)
		}

		args = append(args, reflect.ValueOf(playerId))
		args = append(args, reflect.ValueOf(selectQuadrantMessageObject).Elem())
	case string(ludo_board_constants.MOVE_PAWN):
		pawnMoveMessage := &pawn.PawnMoveMessage{}
		pawnMoveMessageObject, err := pawnMoveMessage.ToObject(string(rawBytes))
		// Log the pawnMoveMessageObject
		// fmt.Println(pawnMoveMessageObject)
		if err != nil {
			log.Printf("Error %s when parsing pawn move message", err)
		}
		args = append(args, reflect.ValueOf(pawnMoveMessageObject).Elem())
	case string(ludo_board_constants.DICEROLL):
		args = append(args, reflect.ValueOf(playerId))
	}

	// log.Printf("Invoking method %s on class %s, with %d arguments", methodName, className, len(args))

	method := reflect.ValueOf(instance).MethodByName(methodName)

	if !method.IsValid() {
		return fmt.Errorf("method %s not found on class %s", methodName, className)
	}

	method.Call(args)

	return nil
}

func (gs *LudoGameService) AddPlayer(boardId string, playerId string, name string, walletAddress string) error {
	// Print roomId, playerId and name
	// fmt.Println(boardId, playerId, name)

	boardInstance, exists := BoardInstances[boardId]

	// log.Println("AddPlayer called with boardId: ", boardId)

	if !exists {
		return fmt.Errorf("game instance not found for room ID: %s", boardId)
	}

	err := boardInstance.AddPlayer(playerId, name, walletAddress)

	if err != nil {
		return fmt.Errorf("error adding player to board: %v", err)
	}

	return nil
}

func (gs *LudoGameService) HandleDisconnection(boardId string, playerId string) error {

	boardInstance, exists := BoardInstances[boardId]

	if !exists {
		// log.Printf("Game instance not found for room ID: %s", boardId)
		return fmt.Errorf("game instance not found for room ID: %s", boardId)
	}

	error := boardInstance.HandleDisconnection(playerId)

	if error != nil {
		log.Printf("Error handling disconnection for player %s in board %s: %v", playerId, boardId, error)
		return error
	}

	return nil
}

func (gs *LudoGameService) CreateTwoPlayerEmptyBoardInstances() error {

	err := gs.createEmptyBoardInstances(2)

	if err != nil {
		log.Printf("Error creating empty board instances: %v", err)
		return err
	}

	return nil
}

func (gs *LudoGameService) CreateFourPlayerEmptyBoardInstances() error {

	err := gs.createEmptyBoardInstances(4)

	if err != nil {
		log.Printf("Error creating empty board instances: %v", err)
		return err
	}

	return nil
}

func (gs *LudoGameService) createEmptyBoardInstances(playerCount int) error {
	// log.Printf("Creating empty board instances for %d player count", playerCount)

	// Count existing waiting boards for this player count
	waitingBoardCount := 0
	existingBoards := make(map[string]*board.Board)

	for id, board := range BoardInstances {
		if board.GetBoardStatus() == ludo_board_constants.WAITING &&
			board.GetMaxPlayers() == playerCount &&
			len(board.GetPlayers()) == 0 {
			waitingBoardCount++
			existingBoards[id] = board
			if waitingBoardCount >= 6 {
				// log.Printf("Already have 6 waiting boards for %d players", playerCount)
				return nil
			}
		}
	}

	// Create only the needed number of boards
	boardsNeeded := 6 - waitingBoardCount
	for i := 0; i < boardsNeeded; i++ {

		boardId := primitive.NewObjectID().Hex()
		rakeAmountType := ludo_board_constants.FIXED
		if len(existingBoards) >= 3 {
			rakeAmountType = ludo_board_constants.PERCENTAGE
		}

		amount := ludo_board_constants.TICKET_AMOUNTS[0]

		newBoardConfig := BoardConfig{
			boardId:        boardId,
			playerCount:    playerCount,
			rakeAmountType: rakeAmountType,
			amount:         amount,
		}

		newBoard := gs.createBoard(newBoardConfig)

		BoardInstances[boardId] = newBoard
		existingBoards[boardId] = newBoard
	}

	// log.Printf("Created %d additional waiting boards for %d players. Total: %d",
	// 	boardsNeeded, playerCount, len(existingBoards))
	return nil
}

func (gs *LudoGameService) createBoard(boardConfig BoardConfig) *board.Board {

	boardId := boardConfig.boardId
	playerCount := boardConfig.playerCount
	rakeAmountType := boardConfig.rakeAmountType
	amount := boardConfig.amount

	newBoard := board.NewBoard(
		boardId,
		playerCount,
		ludo_board_constants.AUTO_PLAY,
		amount,
		int(ludo_board_constants.RAKE_AMOUNT[rakeAmountType]),
		rakeAmountType,
		ludo_board_constants.AUTO_PLAY_TIMER,
	)
	log.Printf("Board created with ID: %s and players: %d", boardId, newBoard.GetMaxPlayers())
	newBoard.SetTicketAmount(amount)

	return newBoard

}

func (gs *LudoGameService) CreateEmptyBoardInstances() error {
	// log.Println("Creating empty board instances")

	// Initialize waiting boards map
	waitingBoards := make(map[int]map[string]*board.Board)
	for _, amount := range ludo_board_constants.TICKET_AMOUNTS {
		waitingBoards[amount] = make(map[string]*board.Board)
	}

	// Count waiting boards by ticket amount
	for id, board := range BoardInstances {
		if board.GetBoardStatus() == ludo_board_constants.WAITING && len(board.GetPlayers()) == 0 {
			ticketAmount := board.GetTicketAmount()
			if _, exists := waitingBoards[ticketAmount]; exists {
				waitingBoards[ticketAmount][id] = board
			}
		}
	}

	// Ensure there are always 6 empty instances for 2-player and 4-player games
	for _, playerCount := range ludo_board_constants.PLAYERS_REQUIRED_TO_START_GAME {
		for _, amount := range ludo_board_constants.TICKET_AMOUNTS {
			for len(waitingBoards[amount]) < 6 {
				boardId := primitive.NewObjectID().Hex()
				rakeAmountType := ludo_board_constants.FIXED

				if len(waitingBoards[amount]) >= 3 {
					rakeAmountType = ludo_board_constants.PERCENTAGE
				}

				newBoard := board.NewBoard(boardId, playerCount, ludo_board_constants.AUTO_PLAY, amount, int(ludo_board_constants.RAKE_AMOUNT[rakeAmountType]), rakeAmountType, ludo_board_constants.AUTO_PLAY_TIMER)
				newBoard.SetTicketAmount(amount)

				BoardInstances[boardId] = newBoard
				waitingBoards[amount][boardId] = newBoard
			}
		}
	}

	// Log total boards and waiting boards by ticket amount
	// _ := len(BoardInstances)
	// log.Printf("Total boards: %d", totalBoards)

	return nil
}

func (s *LudoGameService) cleanupAndCreateBoards() error {
	// log.Println("Cleaning up finished boards and creating new empty boards")

	// Remove finished boards
	for id, board := range BoardInstances {
		if board.GetBoardStatus() == ludo_board_constants.FINISHED && board.GetBoardStatus() == ludo_board_constants.DISCARDED {
			// log.Printf("Removing finished board: %s", id)
			delete(BoardInstances, id)
		}
	}
	// Create new empty boards
	if err := s.CreateTwoPlayerEmptyBoardInstances(); err != nil {
		log.Printf("Error creating 2 player empty board instances: %v", err)
	}

	if err := s.CreateFourPlayerEmptyBoardInstances(); err != nil {
		log.Printf("Error creating 4 player empty board instances: %v", err)
	}

	// log.Println("Finished cleaning up and creating new empty boards")
	return nil
}

func (s *LudoGameService) StartBoardManagement() {

	// Create empty board instances on start
	if err := s.cleanupAndCreateBoards(); err != nil {
		log.Printf("Error in board management: %v", err)
		return
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)

		defer ticker.Stop()
		for range ticker.C {
			// log.Println("Starting board management routine")
			// log.Println("Running board cleanup and creation")
			if err := s.cleanupAndCreateBoards(); err != nil {
				log.Printf("Error in board management: %v", err)
				break
			}
		}
	}()
}

func (gs *LudoGameService) GetBoardList() []*board.Board {

	// log.Println("Getting running board lists")

	var boardLists []*board.Board

	for _, board := range BoardInstances {
		boardLists = append(boardLists, board)
	}

	return boardLists
}
