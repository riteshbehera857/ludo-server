package pawn

import (
	"fmt"
	"ludo/ludo_board_constants"
	"strconv"
)

// MoveResult holds the details of a pawn move
type MoveResult struct {
	InitialPosition  int
	InitialIndex     int
	FinalPosition    int
	FinalIndex       int
	IsAtHome         bool
	ValidationErrors []ValidationError
}

// Pawn represents a game piece that players move around the Ludo board
// Each pawn belongs to a player, has a specific color, and follows a predefined path
type Pawn struct {
	color           string                          // Color of the pawn (red, green, yellow, blue)
	name            string                          // Unique identifier for the pawn (e.g. RED_1)
	currentPosition *int                            // Current location of the pawn on the board
	status          ludo_board_constants.PawnStatus // Current state of the pawn (HOME/PLAYING/FINISHED)
	idlePosition    *int                            // Starting position where the pawn begins and can return
	path            []int                           // Path the pawn follows on the board
}

// NewPawn creates and initializes a new Pawn with the specified attributes
// Parameters:
//   - color: The color of the pawn (red, green, yellow, blue)
//   - name: Unique identifier for the pawn (e.g. RED_1)
//   - homePos: Starting coordinate position on the board
//
// Returns a pointer to the newly created Pawn
func NewPawn(color string, name string, idlePosition *int, path []int) *Pawn {
	return &Pawn{
		color:           color,
		name:            name,
		currentPosition: idlePosition,
		status:          ludo_board_constants.PAWN_IDLE,
		idlePosition:    idlePosition,
		path:            path,
	}
}

// GetPosition returns the current coordinate position of the pawn on the board
// Returns:
//   - Coordinate: The current x,y position of the pawn
func (p *Pawn) GetPosition() *int {
	return p.currentPosition
}

// SetPosition updates the pawn's current position to the specified coordinate
// Parameters:
//   - pos: New coordinate position to move the pawn to
func (p *Pawn) SetPosition(pos *int) {
	p.currentPosition = pos
}

// GetColor returns the color assigned to this pawn
// Returns:
//   - string: The pawn's color (red, green, yellow, blue)
func (p *Pawn) GetColor() string {
	return p.color
}

// GetName returns the unique identifier of the pawn
// Returns:
//   - string: The pawn's name (e.g., "RED_1")
func (p *Pawn) GetName() string {
	return p.name
}

// GetStatus returns the current status of the pawn (HOME/PLAYING/FINISHED)
// Returns:
//   - ludo_board_constants.PawnStatus: Current status of the pawn
func (p *Pawn) GetStatus() ludo_board_constants.PawnStatus {
	return p.status
}

// SetStatus updates the pawn's current status.
// Parameters:
//   - status: The new status to set (HOME, PLAYING, FINISHED)
func (p *Pawn) SetStatus(status ludo_board_constants.PawnStatus) {
	p.status = status
}

// IsAtHome checks if the pawn is in its starting position
// Returns:
//   - bool: true if pawn is at home position, false otherwise
func (p *Pawn) IsIdle() bool {
	return p.status == ludo_board_constants.PAWN_IDLE && p.currentPosition == nil
}

// IsAtFinish checks if the pawn has reached its final destination
// Returns:
//   - bool: true if pawn has finished, false otherwise
func (p *Pawn) IsAtFinish() bool {
	if p.currentPosition == nil {
		return false
	}
	return *p.currentPosition == p.path[len(p.path)-1]
}

// Move attempts to move the pawn by the specified number of steps
// Parameters:
//   - steps: Number of steps to move forward
//
// Returns:
//   - error: Error if move is invalid, nil if successful
func (p *Pawn) Move(steps int) ValidationError {

	// Validate the move
	var validationError ValidationError

	if !p.IsValidMove(steps) {
		validationError = ValidationError{
			Message: fmt.Sprintf("invalid move for pawn %s", p.name),
			CurrentLocation: strconv.Itoa(func() int {
				if p.currentPosition == nil {
					return -1
				}
				return *p.currentPosition
			}()),
		}
		return validationError
	}

	nextPos := p.GetNextPosition(steps)
	p.currentPosition = &nextPos

	if p.status == ludo_board_constants.PAWN_IDLE {
		p.status = ludo_board_constants.PAWN_PLAYING
	}

	return ValidationError{}

}

// IsValidMove checks if moving the specified steps is allowed
func (p *Pawn) IsValidMove(steps int) bool {
	// fmt.Printf("Validating move for pawn %s: steps=%d, currentPosition=%v, status=%v\n", p.name, steps, p.currentPosition, p.status)

	if p.status == ludo_board_constants.PAWN_FINISHED {
		// fmt.Printf("Pawn %s is already finished.\n", p.name)
		return false
	}

	// If the pawn is idle and idlePosition is nil, it can only move if the steps are 6
	if p.status == ludo_board_constants.PAWN_IDLE && p.currentPosition == nil && steps != 6 {
		// fmt.Printf("Pawn %s is idle and cannot move because steps are not 6.\n", p.name)
		return false
	}

	nextPos := p.GetNextPosition(steps)

	// Check if the next position is out of bounds
	if nextPos == -1 {
		// fmt.Printf("Next position for pawn %s after moving %d steps is out of bounds.\n", p.name, steps)
		return false
	}

	isValid := p.isValidPosition(nextPos)
	// fmt.Printf("Is next position %d valid for pawn %s: %v\n", nextPos, p.name, isValid)

	return isValid
}

// GetNextPosition calculates the position after moving specified steps
func (p *Pawn) GetNextPosition(steps int) int {
	currentIndex := p.GetCurrentPathIndex()
	// fmt.Printf("Current index for pawn %s: %d\n", p.name, currentIndex)
	// Log the path
	// fmt.Printf("Path for pawn %s: %v\n", p.name, p.path)

	if currentIndex == -1 {
		if p.currentPosition != nil {
			// fmt.Printf("Pawn %s is not on the path, returning current position: %d\n", p.name, *p.currentPosition)
			return *p.currentPosition
		}
		// fmt.Printf("Pawn %s has an invalid current position\n", p.name)
		if len(p.path) > 0 {
			return p.path[0] // Set to the first index of the path
		}
		// fmt.Printf("Path is empty for pawn %s\n", p.name)
		return -1 // or some default value indicating an invalid position
	}

	nextIndex := currentIndex + steps
	// fmt.Printf("Next index for pawn %s after moving %d steps: %d\n", p.name, steps, nextIndex)

	if nextIndex >= len(p.path) {
		if p.currentPosition != nil {
			// fmt.Printf("Next index %d is out of path bounds for pawn %s, returning current position: %d\n", nextIndex, p.name, *p.currentPosition)
			return -1
		}
		// fmt.Printf("Next index %d is out of path bounds for pawn %s and current position is invalid\n", nextIndex, p.name)
		return -1 // or some default value indicating an invalid position
	}

	// fmt.Printf("Next position for pawn %s: %d\n", p.name, p.path[nextIndex])
	return p.path[nextIndex]
}

// IsBlocked checks if the pawn's next move is blocked by other pawns
// Parameters:
//   - otherPawns: Slice of other pawns to check for collision
//
// Returns:
//   - bool: true if path is blocked, false if clear
func (p *Pawn) IsBlocked(otherPawns []*Pawn) bool {
	nextPos := p.GetNextPosition(1)
	for _, other := range otherPawns {
		if other != p && other.GetPosition() == &nextPos {
			return true
		}
	}
	return false
}

// getCurrentPathIndex finds the current position's index in the path
// Returns:
//   - int: Index in path array, -1 if not found
func (p *Pawn) GetCurrentPathIndex() int {
	if p.currentPosition == nil {
		// fmt.Printf("Pawn %s has no current position set (currentPosition is nil)\n", p.name)
		return -1
	}
	// fmt.Printf("Pawn %s current position: %d\n", p.name, *p.currentPosition)
	for i, pos := range p.path {
		// fmt.Printf("Checking path index %d with position %d for pawn %s\n", i, pos, p.name)
		if pos == *p.currentPosition {
			// fmt.Printf("Pawn %s found at index %d in path\n", p.name, pos)
			return i
		}
	}
	// fmt.Printf("Pawn %s not found in path\n", p.name)
	return -1
}

// isValidPosition checks if a coordinate exists in the pawn's path
// Parameters:
//   - pos: Coordinate to validate
//
// Returns:
//   - bool: true if position is valid, false otherwise
func (p *Pawn) isValidPosition(pos int) bool {
	for _, pathPos := range p.path {
		if pathPos == pos {
			return true
		}
	}
	return false
}

// MovePawn handles the movement of the pawn
// Parameters:
//   - steps: Number of steps to move the pawn
//   - board: The game board
//   - quadrantName: The name of the quadrant whose player is moving the pawn
//
// Returns:
//   - MoveResult: Struct containing the details of the move
//   - error: Error if the move is invalid, nil if successful
func (p *Pawn) MovePawn(steps int, quadrantName string) (MoveResult, error) {

	// Get initial position and index
	initialPosition := -1
	initialIndex := -1
	if p.currentPosition != nil {
		initialPosition = *p.currentPosition
		initialIndex = p.GetCurrentPathIndex()
	}

	// Move the pawn
	validationError := p.Move(steps)

	if validationError.Message != "" {
		return MoveResult{
			IsAtHome:         p.IsIdle(),
			ValidationErrors: []ValidationError{validationError},
		}, nil
	}

	// Get final position and index
	finalPosition := p.path[0] // Set to the first index of the path if initial position is nil
	finalIndex := 0
	if p.currentPosition != nil {
		finalPosition = *p.currentPosition
		finalIndex = p.GetCurrentPathIndex()
	}

	// Check if the pawn has reached the finish
	isAtHome := p.IsAtFinish()

	return MoveResult{
		InitialPosition:  initialPosition,
		InitialIndex:     initialIndex,
		FinalPosition:    finalPosition,
		FinalIndex:       finalIndex,
		IsAtHome:         isAtHome,
		ValidationErrors: []ValidationError{validationError},
	}, nil
}
