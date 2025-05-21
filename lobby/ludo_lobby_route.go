package lobby

import (
	"encoding/json"
	"lobby/response_codes"
	"ludo"
	"ludo/ludo_board_constants"
	"net/http"
	// "ludo"
)

type LudoGameHandler struct{}

type Player struct {
	PlayerId string `json:"playerId"`
	Name     string `json:"name"`
}

type BoardResult struct {
	BoardId                    string   `json:"boardId"`
	Players                    []Player `json:"players"`
	PlayersRequiredToStartGame int      `json:"playersRequiredToStartGame"`
	Status                     string   `json:"status"`
	AutoPlay                   bool     `json:"autoPlay"`
	TicketAmount               int      `json:"ticketAmount"`
}

type BoardListResponse struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Boards  []BoardResult `json:"boards"`
}

func (h *LudoGameHandler) GetBoardList(w http.ResponseWriter, r *http.Request) {

	var response BoardListResponse

	ludo := &ludo.LudoGameService{}

	boardList := ludo.GetBoardList()

	for _, board := range boardList {

		if board.GetBoardStatus() != ludo_board_constants.WAITING && board.GetBoardStatus() != ludo_board_constants.PLAYING {
			continue
		}

		players := []Player{}

		for _, player := range board.GetPlayers() {
			players = append(players, Player{
				PlayerId: player.ID,
				Name:     player.Name,
			})
		}

		response.Boards = append(response.Boards, BoardResult{
			BoardId:                    board.GetID(),
			Players:                    players,
			PlayersRequiredToStartGame: board.GetMaxPlayers(),
			Status:                     string(board.GetBoardStatus()),
			AutoPlay:                   board.GetAutoPlay(),
			TicketAmount:               board.GetTicketAmount(),
		})
	}

	responseCodes := response_codes.GetResponseCodeDetails("BOARD_LIST_FETCHED_SUCCESSFULLY")

	response.Code = responseCodes.Code
	response.Message = responseCodes.Message

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
