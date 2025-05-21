package socket

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("hmXS5XasJL5VSWZ1HenWEV5HvBXSZBdhw8YgKfDNQ+8=")

type AuthClaim struct {
	PlayerId string
	Name     string
	Exp      int64
}

func NewAuthClaim(playerId string, name string) AuthClaim {
	return AuthClaim{
		PlayerId: playerId,
		Name:     name,
	}
}

func VerifyToken(tokenString string) (string, string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	var playerId string
	var name string

	if err != nil {
		return playerId, name, err
	}

	if !token.Valid {
		return playerId, name, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return playerId, name, fmt.Errorf("invalid claims")
	}

	playerIdClaim, ok := claims["playerId"].(string)

	if !ok {
		return playerId, name, fmt.Errorf("invalid playerId")
	}
	playerId = playerIdClaim

	nameClaim, ok := claims["name"].(string)
	if !ok {
		return playerId, name, fmt.Errorf("invalid name")
	}
	name = nameClaim

	exp, ok := claims["exp"].(float64)
	if !ok {
		return playerId, name, fmt.Errorf("invalid exp")
	}

	if int64(exp) < time.Now().Unix() {
		return playerId, name, fmt.Errorf("token expired")
	}

	fmt.Printf("PlayerId: %s, Name: %s\n", playerId, name)

	return playerId, name, nil
}

func CreateToken(playerId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"playerId": playerId,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
			"name":     playerId,
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
