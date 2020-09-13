package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetNextSequence generate new value for a sequence
func GetNextSequence(client *mongo.Client, name string) (int, error) {
	counters := client.Database("deguvon").Collection("counters")
	_, err := counters.UpdateOne(
		context.TODO(),
		bson.M{"_id": name},
		bson.D{{"$inc", bson.D{{"seq", 1}}}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return 0, err
	}
	var sq Serial
	if err = counters.FindOne(context.TODO(), bson.M{"_id": name}).Decode(&sq); err != nil {
		return 0, err
	}
	return sq.Seq, nil
}
