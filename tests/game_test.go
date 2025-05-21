package tests

import (
	"log"
	"metagame/gameserver/game/ludo"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Function to simulate user choosing a quadrant
func chooseQuadrant(availableQuadrants map[string]string, index int) string {
	keys := make([]string, 0, len(availableQuadrants))
	for key := range availableQuadrants {
		keys = append(keys, key)
	}
	return keys[index%len(keys)]
}

// func TestCreateGame(t *testing.T) {
// 	// Create a new game instance
// 	game := ludo.NewGame()

// 	fmt.Print("Game: ", game)

// 	// Verify the initial state of the game
// 	assert.NotNil(t, game, "The game instance should not be nil")
// 	assert.Equal(t, ludo.WAITING, game.GetBoardStatus(), "The initial game status should be WAITING")
// 	assert.Equal(t, 0, len(game.GetPlayers()), "The initial number of players should be 0")
// 	assert.NotNil(t, game.GetBoard(), "The game board should not be nil")
// }

func TestAddPlayers(t *testing.T) {

	quadrantMap := map[string]interface{}{
		"ID": primitive.NewObjectID().Hex,
		"QUADRANT_1": map[string]string{
			"Name":  "QUADRANT_1",
			"Color": "RED",
		},
		"QUADRANT_2": map[string]string{
			"Name":  "QUADRANT_2",
			"Color": "GREEN",
		},
		"QUADRANT_3": map[string]string{
			"Name":  "QUADRANT_3",
			"Color": "YELLOW",
		},
		"QUADRANT_4": map[string]string{
			"Name":  "QUADRANT_4",
			"Color": "BLUE",
		},
	}

	game := ludo.NewGame(quadrantMap)

	// Verify the initial state of the game
	assert.NotNil(t, game, "The game instance should not be nil")
	assert.Equal(t, ludo.WAITING, game.GetBoardStatus(), "The initial game status should be WAITING")
	assert.Equal(t, 0, len(game.GetPlayers()), "The initial number of players should be 0")
	assert.NotNil(t, game.GetBoard(), "The game board should not be nil")
	assert.Equal(t, 4, len(game.GetBoard().GetAvailableQuadrants()), "The number of available quadrants should be 4")

	// players := make([]string, 6)

	for i := 0; i < 4; i++ {
		playerId := primitive.NewObjectID().Hex()

		log.Print("PlayerID: ", playerId)

		availableQuadrants := game.GetBoard().GetAvailableQuadrants()

		// log.Printf("Available Quadrants: %v", availableQuadrants)
		quadrant := chooseQuadrant(availableQuadrants, i)
		// // fmt.Print("Available Quadrants ", availableQuadrants)
		log.Println("Quadrant ", quadrant)

		// Select quadrant for player
		player, err := game.GetBoard().SelectQuadrant(playerId, quadrant)

		if err != nil {
			log.Printf("Error selecting quadrant: %v", err)
		}

		// log.Printf("Player: %v", player)

		// On each quadrantSelection the availableQuadrants value should be less by 1
		assert.Equal(t, 4-(i+1), len(game.GetBoard().GetAvailableQuadrants()), "The number of available quadrants should be 4")

		// If i < 4, and there shouldn't be any quadrant left to choose

		// players[i] = playerID
		// players[i] = ludo.NewPlayer(playerID, quadrant)
	}

	// fmt.Print("Players ", players)
}
