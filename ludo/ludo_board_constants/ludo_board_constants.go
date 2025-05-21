package ludo_board_constants

type BoardStatus string

const (
	WAITING   BoardStatus = "WAITING"   // Waiting for players
	PLAYING   BoardStatus = "PLAYING"   // Game in progress
	FINISHED  BoardStatus = "FINISHED"  // Game completed
	DISCARDED BoardStatus = "DISCARDED" // Game discarded
)

// PawnStatus represents the current state of a pawn in the game
// It can be HOME (not yet started), PLAYING (on board), or FINISHED (reached end)
type PawnStatus string
type InstanceName string
type MethodName string

const (
	// HOME indicates the pawn is in its starting position and hasn't entered the game
	PAWN_IDLE     PawnStatus = "PAWN_IDLE"
	PAWN_PLAYING  PawnStatus = "PAWN_PLAYING"
	PAWN_FINISHED PawnStatus = "PAWN_FINISHED"
)

const (
	GAME_INITIALIZE          = "Game.Initialize"
	SELECT_QUADRANT          = "Select.Quadrant"
	QUADRANT_SELECT          = "Board.SelectQuadrant"
	BOARD_JOINED             = "Board.Joined"
	GAME_START               = "Game.Start"
	TURN                     = "Turn"
	BOARD_DICEROLL           = "Board.DiceRoll"
	BOARD_DICEROLLED         = "Board.DiceRolled"
	BOARD_MOVEPAWN           = "Board.MovePawn"
	BOARD_MOVINGPAWN         = "Board.MovingPawn"
	BOARD_PAWNMOVED          = "Board.PawnMoved"
	BOARD_TURN_COMPLETED     = "Board.TurnCompleted"
	GAME_WINNER              = "Game.Winner"
	GAME_END                 = "Game.End"
	BOARD_DICEROLLING        = "Board.DiceRolling"
	PLAYER_DISCONNECTED      = "Player.Disconnected"
	BOARD_RECONNECTION       = "Player.Reconnected"
	BOARD_BET_FAILED         = "Board.BetFailed"
	BOARD_WAITING_PLAYERS    = "Board.WaitingPlayers"
	BOARD_SELECTING_QUADRANT = "Board.SelectingQuadrant"
)

const (
	RS405 = "Insufficient balance"
)

const (
	DICE     InstanceName = "Dice"
	BOARD    InstanceName = "Board"
	QUADRANT InstanceName = "Quadrant"
)

const (
	SELECT    MethodName = "SelectQuadrant"
	MOVE_PAWN MethodName = "MovePawn"
	DICEROLL  MethodName = "DiceRoll"
)

var SafePositions = []int{91, 36, 23, 102, 133, 188, 201, 122}

var QuadrantsNames = map[int]string{
	1: "QUADRANT_1",
	2: "QUADRANT_2",
	3: "QUADRANT_3",
	4: "QUADRANT_4",
}

var QuadrantsColors = map[string]string{
	"QUADRANT_1": "RED",
	"QUADRANT_2": "GREEN",
	"QUADRANT_3": "YELLOW",
	"QUADRANT_4": "BLUE",
}

var QuadrantsPaths = map[string][]int{
	"QUADRANT_1": {
		91, 92, 93, 94, 95,
		81, 66, 51, 36, 21, 6,
		7,
		8, 23, 38, 53, 68, 83,
		99, 100, 101, 102, 103, 104,
		119,
		134, 133, 132, 131, 130, 129,
		143, 158, 173, 188, 203, 218,
		217,
		216, 201, 186, 171, 156, 141,
		125, 124, 123, 122, 121, 120,
		105,
		106, 107, 108, 109, 110, 111,
	},
	"QUADRANT_2": {
		23, 38, 53, 68, 83,
		99, 100, 101, 102, 103, 104,
		119,
		134, 133, 132, 131, 130, 129,
		143, 158, 173, 188, 203, 218,
		217,
		216, 201, 186, 171, 156, 141,
		125, 124, 123, 122, 121, 120,
		105,
		90, 91, 92, 93, 94, 95,
		81, 66, 51, 36, 21, 6,
		7,
		22, 37, 52, 67, 82, 97,
	},
	"QUADRANT_3": {
		133, 132, 131, 130, 129,
		143, 158, 173, 188, 203, 218,
		217,
		216, 201, 186, 171, 156, 141,
		125, 124, 123, 122, 121, 120,
		105,
		90, 91, 92, 93, 94, 95,
		81, 66, 51, 36, 21, 6,
		7,
		8, 23, 38, 53, 68, 83,
		99, 100, 101, 102, 103, 104,
		119,
		118, 117, 116, 115, 114, 113,
	},
	"QUADRANT_4": {
		201, 186, 171, 156, 141,
		125, 124, 123, 122, 121, 120,
		105,
		90, 91, 92, 93, 94, 95,
		81, 66, 51, 36, 21, 6,
		7,
		8, 23, 38, 53, 68, 83,
		99, 100, 101, 102, 103, 104,
		119,
		134, 133, 132, 131, 130, 129,
		143, 158, 173, 188, 203, 218,
		217,
		202, 187, 172, 157, 142, 127,
	},
}

var PLAYERS_REQUIRED_TO_START_GAME = []int{2, 4}

const (
	AUTO_PLAY = true
)

var NUMBER_OF_BOARDS_REQUIRED_FOR_EACH_PLAYER_COUNT = 6

var NUMBER_OF_BOARDS_REQUIRED_FOR_EACH_AMOUNT = 2
var TICKET_AMOUNTS = []int{100, 200, 500}

var RAKE_AMOUNT = map[RakeAmountType]float64{
	FIXED:      0.0,
	PERCENTAGE: 10.0,
}

type RakeAmountType string

var AUTO_PLAY_TIMER = 5

const (
	FIXED      RakeAmountType = "FIXED"
	PERCENTAGE RakeAmountType = "PERCENTAGE"
)
