package api

import "net/http"

func HandleRoutes() {
	boardRouteHandler := &BoardRouteHandler{}

	http.HandleFunc("/api/board/list", boardRouteHandler.GetBoardList)
}
