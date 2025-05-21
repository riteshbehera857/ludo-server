package board

import (
	"ludo/ludo_board_constants"
	"ludo/player"
	"time"
)

// Move represents a single move made by a pawn
type MoveSchema struct {
	InitialPosition int       `bson:"initialPosition" json:"initialPosition"`
	FinalPosition   int       `bson:"finalPosition" json:"finalPosition"`
	Timestamp       time.Time `bson:"timestamp" json:"timestamp"`
}

// Game represents the overall game state
type BoardSchema struct {
	ID                         string                              `bson:"_id" json:"_id"`
	BoardId                    string                              `bson:"boardId" json:"boardId"`
	TicketAmount               int                                 `bson:"ticketAmount" json:"ticketAmount"`
	RakeAmount                 int                                 `bson:"rakeAmount" json:"rakeAmount"`
	RakeAmountType             ludo_board_constants.RakeAmountType `bson:"rakeAmountType" json:"rakeAmountType"`
	WinningAmount              int                                 `bson:"winningAmount" json:"winningAmount"`
	Status                     ludo_board_constants.BoardStatus    `bson:"status" json:"status"`
	AutoPlay                   bool                                `bson:"autoPlay" json:"autoPlay"`
	AutoPlayTimer              int                                 `bson:"autoPlayTimer" json:"autoPlayTimer"`
	PlayersRequiredToStartGame int                                 `bson:"playersRequiredToStartGame" json:"playersRequiredToStartGame"`
	StartTime                  *time.Time                          `bson:"startTime,omitempty" json:"startTime,omitempty"`
	EndTime                    *time.Time                          `bson:"endTime,omitempty" json:"endTime,omitempty"`
	Winner                     *string                             `bson:"winner,omitempty" json:"winner,omitempty"`
	Players                    []player.PlayerSchema               `bson:"players" json:"players"`
	PawnMoves                  map[string]map[string][]MoveSchema  `bson:"pawnMoves" json:"pawnMoves"`
}
