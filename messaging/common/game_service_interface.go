package common

type GameService interface {
	ProcessMessage(boardId string, playerId string, message Message, rawBytes []byte) error
	AddPlayer(boardId string, playerId string, name string, walletAddress string) error
	HandleDisconnection(boardId string, playerId string) error
	CreateEmptyBoardInstances() error
	StartBoardManagement()
}
