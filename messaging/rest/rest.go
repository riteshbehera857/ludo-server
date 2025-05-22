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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Lobby server is live!")
	})

	log.Printf("Starting lobby server on port %s", port)
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
