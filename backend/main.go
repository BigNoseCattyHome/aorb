package main

import (
	"context"
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 运行整个项目

func main() {
	collection := database.MongoDbClient.Database("aorb").Collection("polls")

	// Whether user had already voted or not
	filter4Check := bson.D{
		{"pollUuid", "dfbeb25d-f79d-4003-aafb-995ea6fe3453"},
		{"voteList.voteUserName", "aaa"},
	}
	projection := bson.D{
		{"voteList.voteUserName", 1},
	}

	var result bson.D
	collection.FindOne(context.TODO(), filter4Check, options.FindOne().SetProjection(projection)).Decode(&result)

	fmt.Println(result != nil)
}
