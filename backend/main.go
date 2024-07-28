package main

import (
	"context"
	"fmt"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/go-services/poll/models"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/cached"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"go.mongodb.org/mongo-driver/bson"
)

// 运行整个项目

func main() {

	var tempPoll pollModels.Poll
	_, err := cached.GetWithFunc(context.TODO(), fmt.Sprintf("PollExistedCached-%s", "8d01d19c-9d6e-41e9-ae1c-5060076de686"), func(ctx context.Context, key string) (string, error) {
		collection := database.MongoDbClient.Database("aorb").Collection("polls")
		cursor := collection.FindOne(ctx, bson.M{"pollUuid": "8d01d19c-9d6e-41e9-ae1c-5060076de686"})
		if cursor.Err() != nil {
			return "false", cursor.Err()
		}
		if err := cursor.Decode(&tempPoll); err != nil {
			return "false", err
		}
		return "true", nil
	})

	fmt.Println(tempPoll)
	fmt.Println(err)

}
