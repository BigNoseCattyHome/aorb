package main

import (
	"context"
	"fmt"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/go-services/poll/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/cached"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"go.mongodb.org/mongo-driver/bson"
)

// 运行整个项目

func main() {

	var req *poll.PollExistRequest = &poll.PollExistRequest{}
	req.PollId = 1
	var tempPoll pollModels.Poll
	_, err := cached.GetWithFunc(context.Background(), fmt.Sprintf("PollExistedCached-%d", req.PollId), func(ctx context.Context, key string) (string, error) {
		collection := database.MongoDbClient.Database("aorb").Collection("polls")
		cursor := collection.FindOne(ctx, bson.M{"_id": req.PollId})
		if cursor.Err() != nil {
			return "false", cursor.Err()
		}
		if err := cursor.Decode(&tempPoll); err != nil {
			return "false", err
		}
		fmt.Println(tempPoll)
		return "true", nil
	})
	fmt.Println(err)
}
