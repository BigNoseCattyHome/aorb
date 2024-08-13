package main

import (
	"context"
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"go.mongodb.org/mongo-driver/bson"
)

// 运行整个项目

func main() {

	userCollection := database.MongoDbClient.Database("aorb").Collection("users")

	// 查看是否已经关注
	filter4Check := bson.D{{"username", "user1"}}
	cursor := userCollection.FindOne(context.TODO(), filter4Check)
	var result bson.M
	cursor.Decode(&result)

	fmt.Println(result["followed"].(bson.M)["usernames"].(bson.A))
}
