package board

import (
	"context"
	"fmt"
	"log"
	"ludo/ludo_board_constants"
	"ludo/player"
	"metagame/gameserver/config"
	"metagame/gameserver/helpers"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardDAO struct {
	collection *mongo.Collection
}

func NewBoardDAO() *BoardDAO {
	client := helpers.GetMongoClient()
	collection := client.Database(config.GetConfig().Database).Collection("ludo_games")
	// log.Println("NewBoardDAO: Initialized BoardDAO")
	return &BoardDAO{
		collection: collection,
	}
}

func (dao *BoardDAO) UpdatePlayerConnectionDetails(boardId string, playerId string, cType string, time time.Time) error {

	field := "disconnectedAt"

	if cType == "reconnection" {
		field = "reconnectedAt"
	}

	// log.Printf("UpdatePlayerDisconnectionTime: Updating disconnection time for player %s in boardId %s", playerId, boardId)

	filter := bson.M{"boardId": boardId, "players.playerId": playerId}

	update := bson.M{"$set": bson.M{fmt.Sprintf("players.$.%s", field): time}}

	_, err := dao.collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Printf("UpdatePlayerDisconnectionTime: Error updating disconnection time for player %s in boardId %s: %v", playerId, boardId, err)
		return fmt.Errorf("failed to update disconnection time for player %s in game with ID %s: %v", playerId, boardId, err)
	}

	// log.Printf("UpdatePlayerDisconnectionTime: Successfully updated disconnection time for player %s in boardId %s", playerId, boardId)

	return nil
}

func (dao *BoardDAO) InsertBoard(board BoardSchema) (primitive.M, error) {
	// log.Println("InsertBoard: Inserting new board")
	result, err := dao.collection.InsertOne(context.Background(), board)
	if err != nil {
		log.Printf("InsertBoard: Error inserting board: %v", err)
		return nil, err
	}

	// Find the inserted document
	var insertedBoard BoardSchema
	err = dao.collection.FindOne(context.Background(), bson.M{"_id": result.InsertedID}).Decode(&insertedBoard)
	if err != nil {
		return nil, err
	}

	// log.Printf("InsertBoard: Board inserted with ID %v", result.InsertedID)
	return primitive.M{"boardId": insertedBoard.BoardId}, nil
}

func (dao *BoardDAO) UpdateGameStatusAndAddStartTime(boardId string, status ludo_board_constants.BoardStatus, sTime time.Time) error {
	// log.Printf("UpdateGameStatusAndAddStartTime: Updating status to %v and start time for boardId %s", status, boardId)
	filter := bson.M{"boardId": boardId}

	update := bson.M{"$set": bson.M{"status": status, "startTime": sTime}}

	_, err := dao.collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Printf("UpdateGameStatusAndAddStartTime: Error updating status for boardId %s: %v", boardId, err)
		return fmt.Errorf("failed to update status for game with ID %s: %v", boardId, err)
	}

	// log.Printf("UpdateGameStatusAndAddStartTime: Successfully updated status for boardId %s", boardId)
	return nil
}

func formatQuadrantAndPawnNames(quadrant, pawn string) (string, string) {
	quadrantParts := strings.Split(quadrant, "_")

	pawnParts := strings.Split(pawn, "_")

	formattedQuadrant := fmt.Sprintf("quadrant-%s", strings.ToLower(quadrantParts[1]))

	formattedPawn := fmt.Sprintf("pawn-%s", strings.ToLower(pawnParts[3]))

	return formattedQuadrant, formattedPawn
}

func (dao *BoardDAO) UpdateBoardPawnMovement(boardId, quadrant, pawn string, initialPosition, finalPosition int, timestamp time.Time, diceResult int) error {
	// log.Printf("UpdateBoardPawnMovement: Updating pawn movement for boardId %s, quadrant %s, pawn %s", boardId, quadrant, pawn)
	formattedQuadrant, formattedPawn := formatQuadrantAndPawnNames(quadrant, pawn)

	filter := bson.M{"boardId": boardId}
	update := bson.M{
		"$push": bson.M{
			fmt.Sprintf("pawnMoves.%s.%s", formattedQuadrant, formattedPawn): bson.M{
				"diceResult":      diceResult,
				"initialPosition": initialPosition,
				"finalPosition":   finalPosition,
				"timestamp":       timestamp,
			},
		},
	}

	_, err := dao.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("UpdateBoardPawnMovement: Error updating pawn movement for boardId %s: %v", boardId, err)
		return fmt.Errorf("failed to update pawn movement for game with ID %s: %v", boardId, err)
	}

	// log.Printf("UpdateBoardPawnMovement: Successfully updated pawn movement for boardId %s", boardId)
	return nil
}

func (dao *BoardDAO) GetBoardById(boardId string) (*BoardSchema, error) {

	// log.Println("GetBoardById: Fetching board with boardId: ", boardId)

	var board BoardSchema

	filter := bson.M{"boardId": boardId}

	err := dao.collection.FindOne(context.Background(), filter).Decode(&board)

	if err != nil {
		log.Printf("GetBoardById: Error fetching board with boardId %s: %v", boardId, err)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("game with ID %s not found", boardId)
		}
		return nil, err
	}
	// log.Printf("GetBoardById: Successfully fetched board with boardId %s", boardId)
	return &board, nil
}

func (dao *BoardDAO) AddPlayerToBoard(boardId string, player player.PlayerSchema) error {
	// log.Printf("AddPlayerToBoard: Adding player to boardId %s", boardId)
	filter := bson.M{"boardId": boardId}

	update := bson.M{"$push": bson.M{"players": player}}

	_, err := dao.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("AddPlayerToBoard: Error adding player to boardId %s: %v", boardId, err)
		return fmt.Errorf("failed to add player to game with ID %s: %v", boardId, err)
	}

	// log.Printf("AddPlayerToBoard: Successfully added player to boardId %s", boardId)
	return nil
}

func (dao *BoardDAO) UpdateBoardStatusAndAddEndTime(boardId string, status ludo_board_constants.BoardStatus, eTime time.Time) error {
	// log.Printf("UpdateBoardStatusAndAddEndTime: Updating status to %v and end time for boardId %s", status, boardId)
	filter := bson.M{"boardId": boardId}

	update := bson.M{"$set": bson.M{"status": status, "endTime": eTime}}

	_, err := dao.collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Printf("UpdateBoardStatusAndAddEndTime: Error updating status for boardId %s: %v", boardId, err)
		return fmt.Errorf("failed to update status for game with ID %s: %v", boardId, err)
	}

	// log.Printf("UpdateBoardStatusAndAddEndTime: Successfully updated status for boardId %s", boardId)
	return nil

}

func (dao *BoardDAO) UpdateGameWinner(boardId, winner string, winningAmount int) error {

	// log.Printf("UpdateGameWinner: Updating winner to %s with winning amount %d for boardId %s", winner, winningAmount, boardId)
	filter := bson.M{"boardId": boardId}

	update := bson.M{"$set": bson.M{"winner": winner, "winningAmount": winningAmount}}

	_, err := dao.collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Printf("UpdateGameWinner: Error updating winner for boardId %s: %v", boardId, err)
		return fmt.Errorf("failed to update winner for game with ID %s: %v", boardId, err)
	}

	// log.Printf("UpdateGameWinner: Successfully updated winner for boardId %s", boardId)
	return nil

}
