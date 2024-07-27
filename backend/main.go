package main

import (
	"context"
	"fmt"
	pollModels "github.com/BigNoseCattyHome/aorb/backend/go-services/poll/models"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

// 运行整个项目

func main() {

	var result []pollModels.Poll
	collections := database.MongoDbClient.Database("aorb").Collection("polls")
	filter := bson.D{}
	cursor, _ := collections.Find(context.TODO(), filter)
	cursor.All(context.TODO(), &result)
	fmt.Println(reflect.TypeOf(result))
	fmt.Println(result[0].CommentList)
}
