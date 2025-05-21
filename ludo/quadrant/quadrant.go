package quadrant

import (
	"errors"
	"fmt"
	"log"

	"ludo/pawn"
	"ludo/player"

	"go.mongodb.org/mongo-driver/mongo"
)

type Quadrant struct {
	name       string         // Name of the quadrant (QUADRANT_1/QUADRANT_2/QUADRANT_3/QUADRANT_4)
	color      string         // Color of the quadrant (red/blue/green/yellow)
	path       []int          // Path pawns follow from this quadrant
	pawns      []*pawn.Pawn   // Pawns currently in this quadrant
	player     *player.Player // Player associated with this quadrant
	isOccupied bool
}

// NewQuadrant creates a new quadrant with the specified color and initializes its pawns
func NewQuadrant(color string, player *player.Player, quadrantName string, path []int) *Quadrant {
	q := &Quadrant{
		name:       quadrantName,
		color:      color,
		pawns:      make([]*pawn.Pawn, 0, 4),
		path:       path,
		player:     player,
		isOccupied: false,
	}

	// Initialize 4 pawns for this quadrant
	for i := 1; i <= 4; i++ {
		pawnName := fmt.Sprintf("%s_%s_%d", quadrantName, "PAWN", i)
		pawn := pawn.NewPawn(color, pawnName, nil, path)
		q.pawns = append(q.pawns, pawn)
	}

	return q
}

// Get Quadrant Name
func (q *Quadrant) GetName() string {
	return q.name
}

// GetPawns returns all pawns in the quadrant
func (q *Quadrant) GetPawns() []*pawn.Pawn {
	if q == nil {
		// log.Println("Quadrant is nil")
		return nil
	}
	if q.pawns == nil {
		// log.Println("Quadrant pawns are nil")
		return nil
	}
	return q.pawns
}

func (q *Quadrant) GetPawnNames() []string {
	names := []string{}
	for _, pawn := range q.pawns {
		names = append(names, pawn.GetName())
	}
	return names
}

func (q *Quadrant) GetIfQuadrantIsOccupied() bool {
	return q.isOccupied
}

func (q *Quadrant) RemovePlayer() {
	if q != nil {
		q.player = nil
		q.SetIsOccupied(false)
	}
}

func (q *Quadrant) SetIsOccupied(isOccupied bool) {
	if q != nil {
		q.isOccupied = isOccupied
	}
}

// Get Path
func (q *Quadrant) GetPath() []int {
	return q.path
}

func (q *Quadrant) GetPlayer() *player.Player {
	return q.player
}

func (q *Quadrant) SetPlayer(player *player.Player) {
	q.player = player
}

// SelectQuadrant allows a player to select a quadrant from the available quadrants and creates a new player
func (q *Quadrant) Select(player *player.Player) (*player.Player, error) {

	if q.player != nil {
		return nil, errors.New("Quadrant already selected")
	}

	q.SetPlayer(player)

	q.SetIsOccupied(true)

	player.AssignQuadrant(q.GetName())

	return player, nil
}

// GetPawnByNumber returns a specific pawn by its number
// Parameters:
//   - number: Pawn number (1-4)
//
// Returns pointer to Pawn if found, nil otherwise
func (q *Quadrant) GetPawnByName(name string) *pawn.Pawn {
	for _, pawn := range q.pawns {
		if pawn.GetName() == name {
			return pawn
		}
	}
	return nil
}

// GetColor returns the color of this quadrant
// Returns quadrant color as string
func (q *Quadrant) GetColor() string {
	return q.color
}

// CountFinishedPawns returns number of pawns that have reached finish
// Returns count of pawns with FINISHED status
func (q *Quadrant) CountFinishedPawns() int {
	count := 0
	for _, pawn := range q.pawns {
		if pawn.IsAtFinish() {
			count++
		}
	}
	return count
}

// HasWon checks if all pawns from this quadrant have finished
// Returns true if all pawns are FINISHED
func (q *Quadrant) HasWon() bool {
	return q.CountFinishedPawns() == 4
}

// GetQuadrantConfig retrieves the quadrant configuration from the database
func GetQuadrantConfig() (string, error) {
	quadrantDao := NewQuadrantDao()

	quadrantConfig, err := quadrantDao.FetchAndProcessQuadrantConfig()

	if err != nil {
		return "", err
	}

	return quadrantConfig, nil
}

func InitializeQuadrantConfigIfNotExists(quadrantNames map[int]string,
	quadrantPaths map[string][]int,
	quadrantColors map[string]string,
	safePositions []int) (string, error) {

	quadrantDao := NewQuadrantDao()

	newQuadrantConfig := NewQuadrantConfig(quadrantNames, quadrantPaths, quadrantColors, safePositions)

	quadrantConfig, err := GetQuadrantConfig()

	if err == mongo.ErrNoDocuments {
		quadrantConfigInsertError := quadrantDao.InsertQuadrantConfig(newQuadrantConfig)
		if quadrantConfigInsertError != nil {
			log.Fatal(quadrantConfigInsertError)
			return "", quadrantConfigInsertError
		}

		quadrantConfig, err = GetQuadrantConfig()

		if err != nil {
			log.Fatal(err)
			return "", err
		}
	}

	return quadrantConfig, nil
}
