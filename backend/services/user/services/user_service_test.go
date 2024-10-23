package services

import (
	"testing"

	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 测试 GetUserInfo 函数
func TestGetUserInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := NewMockCollection(ctrl)
	mockCursor := NewMockCursor(ctrl)

	// 替换 MongoDB 连接
	database.MongoDbClient = NewMockMongoClient(ctrl)
	database.MongoDbClient.EXPECT().Database("aorb").Return(mockCollection)

	username := "testuser"
	fields := []string{"username", "avatar"}

	// 模拟数据库查询返回的用户信息
	expectedUser := user.User{
		Username: username,
		Avatar:   "http://example.com/avatar.png",
	}
	mockCollection.EXPECT().
		FindOne(gomock.Any(), bson.M{"username": username}, gomock.Any()).
		Return(mockCursor)
	mockCursor.EXPECT().Decode(gomock.Any()).SetArg(0, expectedUser).Return(nil)

	resp, err := services.GetUserInfo(username, fields)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", resp.User.Username)
	assert.Equal(t, "http://example.com/avatar.png", resp.User.Avatar)
}

// 测试用户不存在的情况
func TestGetUserInfo_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := NewMockCollection(ctrl)
	mockCursor := NewMockCursor(ctrl)

	database.MongoDbClient = NewMockMongoClient(ctrl)
	database.MongoDbClient.EXPECT().Database("aorb").Return(mockCollection)

	username := "unknown"
	fields := []string{"username", "avatar"}

	// 模拟用户不存在
	mockCollection.EXPECT().
		FindOne(gomock.Any(), bson.M{"username": username}, gomock.Any()).
		Return(mockCursor)
	mockCursor.EXPECT().Decode(gomock.Any()).Return(mongo.ErrNoDocuments)

	resp, err := services.GetUserInfo(username, fields)
	assert.NoError(t, err)
	assert.Equal(t, strings.UnableToQueryUserErrorCode, resp.StatusCode)
}

// 测试 CheckUserExists 函数
func TestCheckUserExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := NewMockCollection(ctrl)
	database.MongoDbClient = NewMockMongoClient(ctrl)
	database.MongoDbClient.EXPECT().Database("aorb").Return(mockCollection)

	// 用户存在的情况
	username := "testuser"
	mockCollection.EXPECT().FindOne(gomock.Any(), bson.M{"username": username}).Return(nil)
	resp, err := services.CheckUserExists(username)
	assert.NoError(t, err)
	assert.True(t, resp.Existed)

	// 用户不存在的情况
	username = "unknown"
	mockCollection.EXPECT().FindOne(gomock.Any(), bson.M{"username": username}).Return(mongo.ErrNoDocuments)
	resp, err = services.CheckUserExists(username)
	assert.NoError(t, err)
	assert.False(t, resp.Existed)
}

// 测试 IsUserFollowing 函数
func TestIsUserFollowing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCollection := NewMockCollection(ctrl)
	mockCursor := NewMockCursor(ctrl)

	database.MongoDbClient = NewMockMongoClient(ctrl)
	database.MongoDbClient.EXPECT().Database("aorb").Return(mockCollection)

	username := "follower"
	targetUsername := "followee"

	// 模拟查询返回的用户信息
	mockCollection.EXPECT().
		FindOne(gomock.Any(), bson.M{"username": username}).
		Return(mockCursor)
	mockCursor.EXPECT().Decode(gomock.Any()).SetArg(0, user.User{
		Followed: &user.FollowedList{Usernames: []string{"followee"}},
	}).Return(nil)

	resp, err := services.IsUserFollowing(username, targetUsername)
	assert.NoError(t, err)
	assert.True(t, resp.IsFollowing)
}
