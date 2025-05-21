package player

import "time"

type PlayerSchema struct {
	ID             string    `bson:"_id" json:"_id"`
	PlayerID       string    `bson:"playerId" json:"playerId"`
	Name           string    `bson:"name" json:"name"`
	Quadrant       string    `bson:"quadrant" json:"quadrant"`
	JoinedAt       time.Time `bson:"joinedAt" json:"joinedAt"`
	DisconnectedAt time.Time `bson:"disconnectedAt,omitempty" json:"disconnectedAt,omitempty"`
	ReconnectedAt  time.Time `bson:"reconnectedAt,omitempty" json:"reconnectedAt,omitempty"`
}
