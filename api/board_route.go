package api

import (
	"encoding/json"
	"net/http"
)

type BoardRouteHandler struct{}

type BoardListResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (brh *BoardRouteHandler) GetBoardList(w http.ResponseWriter, r *http.Request) {
	response := BoardListResponse{
		Status:  "success",
		Message: "Game is running",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
