package player

// Player represents a participant in the Ludo game
type Player struct {
	ID                      string // Unique identifier for the player
	PlayerId                string // Player's ID from platform
	Name                    string // Player's name
	Quadrant                string // Player's color (red, green, yellow, blue)
	QuadrantSelectionStatus int    // 0: Not selected, 1: Sent selection, 2: Selected
	ConnectionStatus        int    // 0: Disconnected, 1: Connected
	BetId                   string // Bet ID for the player
	WalletAddress           string // Wallet address for the player
}

const (
	PLAYER_DISCONNECTED = iota
	PLAYER_CONNECTED
)

// NewPlayer creates a new player with specified ID and color
func NewPlayer(id string, name string, quadrant string, connectionStatus int, walletAddress string) *Player {
	return &Player{
		PlayerId:         id,
		Name:             name,
		Quadrant:         quadrant,
		ConnectionStatus: connectionStatus,
		WalletAddress:    walletAddress,
	}
}

func (p *Player) GetPlayerId() string {
	return p.PlayerId
}

func (p *Player) AssignQuadrant(quadrantName string) {
	p.Quadrant = quadrantName
}

func (p *Player) GetQuadrant() string {
	return p.Quadrant
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) SetBetId(betId string) {
	p.BetId = betId
}

func (p *Player) HasSelectedQuadrant() bool {
	return p.QuadrantSelectionStatus == 2
}

func (p *Player) SetConnectionStatus(status int) {
	p.ConnectionStatus = status
}

func (p *Player) IsConnected() bool {
	return p.ConnectionStatus == PLAYER_CONNECTED
}
