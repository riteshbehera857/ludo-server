package rest

import (
	"fmt"
	"lobby"
	"log"
	"net/http"
)

func StartRESTApiServer(port string) error {
	lobbyRouteHandler := &lobby.LobbyRouteHandler{}

	lobbyRouteHandler.HandleRoutes()

	log.Printf("Starting lobby server on port %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
