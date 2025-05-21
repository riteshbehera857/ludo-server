package quadrant

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"metagame/gameserver/config"
	"metagame/gameserver/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type QuadrantDao struct {
	collection *mongo.Collection
}

func NewQuadrantDao() *QuadrantDao {
	client := helpers.GetMongoClient()
	collection := client.Database(config.GetConfig().Database).Collection("quadrantConfig")
	return &QuadrantDao{
		collection: collection,
	}
}

func (dao *QuadrantDao) GetQuadrantConfig() (QuadrantConfigSchema, error) {
	var quadrantConfig QuadrantConfigSchema
	err := dao.collection.FindOne(context.Background(), bson.M{}).Decode(&quadrantConfig)
	if err != nil {
		return QuadrantConfigSchema{}, err
	}
	return quadrantConfig, nil
}

func (dao *QuadrantDao) InsertQuadrantConfig(quadrantConfig QuadrantConfigSchema) error {
	// log.Println("Inserting QuadrantConfig")
	_, err := dao.collection.InsertOne(context.Background(), quadrantConfig)
	if err != nil {
		return err
	}
	return nil
}

func (dao *QuadrantDao) FetchAndProcessQuadrantConfig() (string, error) {
	if dao.collection == nil {
		return "", fmt.Errorf("QuadrantConfigCollection is not initialized")
	}

	var quadrantConfig primitive.M

	err := dao.collection.FindOne(context.Background(), bson.D{{}}).Decode(&quadrantConfig)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", err
		}
		log.Printf("Error fetching QuadrantConfig: %v", err)
		return "", err
	}

	// Convert _id to ObjectId() format
	if id, ok := quadrantConfig["_id"].(primitive.ObjectID); ok {
		quadrantConfig["_id"] = "ObjectId(" + id.Hex() + ")"
	}

	jsonData, err := json.Marshal(quadrantConfig)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
