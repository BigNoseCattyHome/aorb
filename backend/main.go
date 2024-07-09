package main

import (
	"context"
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/comment/models"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"go.mongodb.org/mongo-driver/bson"
)

// 运行整个项目

func main() {
	collection := database.MongoDbClient.Database("aorb").Collection("comments")
	cursor, _ := collection.Find(context.TODO(), bson.D{})

	var comments []models.Comment
	if err := cursor.All(context.TODO(), &comments); err != nil {
		panic(err)
	}
	for _, comment := range comments {
		fmt.Println(comment)
	}
}
