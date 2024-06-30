package database

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	logging "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDbClient *mongo.Client

func init() {
	logging.Info("Connecting to mongodb")
	clientOptions := options.Client().ApplyURI("mongodb://" + config.Conf.MongoDB.Host + ":" + config.Conf.MongoDB.Port)
	var err error
	MongoDbClient, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logging.Info(err)
	}
	err = MongoDbClient.Ping(context.TODO(), nil)
	if err != nil {
		logging.Info(err)
	}
	logging.Info("mongodb connect success")
}
