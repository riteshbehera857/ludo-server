package quadrant

import (
	// "messaging/common/constants"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuadrantSchema struct {
	Name  string `bson:"name"`
	Path  []int
	Color string `bson:"color"`
}

type QuadrantConfigSchema struct {
	SafePositions []int          `json:"safePositions" bson:"safePositions"`
	QUADRANT_1    QuadrantSchema `json:"QUADRANT_1" bson:""`
	QUADRANT_2    QuadrantSchema `json:"QUADRANT_2" bson:""`
	QUADRANT_3    QuadrantSchema `json:"QUADRANT_3" bson:""`
	QUADRANT_4    QuadrantSchema `json:"QUADRANT_4," bson:""`
}

// NewQuadrantConfig creates a new QuadrantConfig instance.
// Accept arguements for all the quadrants and safe positions.
func NewQuadrantConfig(
	quadrantNames map[int]string,
	quadrantPaths map[string][]int,
	quadrantColors map[string]string,
	safePositions []int,
) QuadrantConfigSchema {
	return QuadrantConfigSchema{
		SafePositions: safePositions,
		QUADRANT_1: QuadrantSchema{
			Path:  quadrantPaths[quadrantNames[1]],
			Name:  quadrantNames[1],
			Color: quadrantColors[quadrantNames[1]],
		},
		QUADRANT_2: QuadrantSchema{
			Path:  quadrantPaths[quadrantNames[2]],
			Name:  quadrantNames[2],
			Color: quadrantColors[quadrantNames[2]],
		},
		QUADRANT_3: QuadrantSchema{
			Path:  quadrantPaths[quadrantNames[3]],
			Name:  quadrantNames[3],
			Color: quadrantColors[quadrantNames[3]],
		},
		QUADRANT_4: QuadrantSchema{
			Path:  quadrantPaths[quadrantNames[4]],
			Name:  quadrantNames[4],
			Color: quadrantColors[quadrantNames[4]],
		},
	}
}

func ConvertToQuadrantMap(quadrantConfig primitive.M) (map[string]interface{}, error) {
	// Convert _id to ObjectId() format
	if id, ok := quadrantConfig["_id"].(primitive.ObjectID); ok {
		quadrantConfig["_id"] = "ObjectId(" + id.Hex() + ")"
	}

	quadrantMap := map[string]interface{}{
		"ID":            quadrantConfig["_id"],
		"SafePositions": convertToIntSlice(quadrantConfig["safePositions"]),
		"QUADRANT_1": map[string]interface{}{
			"Path":  convertToIntSlice(quadrantConfig["quadrant_1"].(map[string]interface{})["path"]),
			"Name":  quadrantConfig["quadrant_1"].(map[string]interface{})["name"].(string),
			"Color": quadrantConfig["quadrant_1"].(map[string]interface{})["color"].(string),
		},
		"QUADRANT_2": map[string]interface{}{
			"Path":  convertToIntSlice(quadrantConfig["quadrant_2"].(map[string]interface{})["path"]),
			"Name":  quadrantConfig["quadrant_2"].(map[string]interface{})["name"].(string),
			"Color": quadrantConfig["quadrant_2"].(map[string]interface{})["color"].(string),
		},
		"QUADRANT_3": map[string]interface{}{
			"Path":  convertToIntSlice(quadrantConfig["quadrant_3"].(map[string]interface{})["path"]),
			"Name":  quadrantConfig["quadrant_3"].(map[string]interface{})["name"].(string),
			"Color": quadrantConfig["quadrant_3"].(map[string]interface{})["color"].(string),
		},
		"QUADRANT_4": map[string]interface{}{
			"Path":  convertToIntSlice(quadrantConfig["quadrant_4"].(map[string]interface{})["path"]),
			"Name":  quadrantConfig["quadrant_4"].(map[string]interface{})["name"].(string),
			"Color": quadrantConfig["quadrant_4"].(map[string]interface{})["color"].(string),
		},
	}

	return quadrantMap, nil
}

func convertToIntSlice(data interface{}) []int {
	if data == nil {
		return nil
	}

	interfaceSlice, ok := data.([]interface{})
	if !ok {
		return nil
	}

	intSlice := make([]int, len(interfaceSlice))
	for i, v := range interfaceSlice {
		switch num := v.(type) {
		case float64:
			intSlice[i] = int(num)
		case int:
			intSlice[i] = num
		case int32:
			intSlice[i] = int(num)
		case int64:
			intSlice[i] = int(num)
		default:
			return nil
		}
	}

	return intSlice
}
