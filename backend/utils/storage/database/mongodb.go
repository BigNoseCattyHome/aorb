package database

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoDbClient = InitMongoDbClient()
	mongoDatabase = mongoDbClient.Database("aorb")
)

func InitMongoDbClient() *mongo.Client {
	logging.Info("Connecting to mongodb")
	clientOptions := options.Client().ApplyURI("mongodb://" + config.Conf.MongoDB.Host + ":" + config.Conf.MongoDB.Port)
	var err error
	mongoDbClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logging.Info(err)
	}
	err = mongoDbClient.Ping(context.TODO(), nil)
	if err != nil {
		logging.Info(err)
	}
	logging.Info("mongodb connect success")
	return mongoDbClient
}
