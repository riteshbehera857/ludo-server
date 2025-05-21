package lobby

import (
	"net/http"
)

type LobbyRouteHandler struct {
}

func (s *LobbyRouteHandler) HandleRoutes() {
	ludoGameHandler := &LudoGameHandler{}

	http.HandleFunc("/api/ludo/board-list", ludoGameHandler.GetBoardList)
}
