package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/conf"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 使用logging库，添加字段日志 UserRpcServerName
var log = logging.LogService(config.UserRpcServerName)

func GetUserInfo(username string, fields []string) (resp *user.UserResponse, err error) {
	// 根据 username 和 fields 获取用户信息
	collection := database.MongoDbClient.Database("aorb").Collection("users")
	filter := bson.M{"username": username}
	projection := bson.M{} // 设置查询的字段
	for _, field := range fields {
		projection[field] = 1 // 1表示要查询的字段, 0表示不查询
	}
	opts := options.FindOne().SetProjection(projection)
	var queryUser user.User
	err = collection.FindOne(context.TODO(), filter, opts).Decode(&queryUser)
	log.Debug("get the user with fields: ", &queryUser)
	if err != nil {
		// 如果查询不到用户，返回 UnableToQueryUserErrorCode "无法查询到对应用户"
		// 但是没有返回错误，只在返回的用户信息中 StatusCode 和 StatusMsg 字段中返回错误信息
		if err == mongo.ErrNoDocuments {
			return &user.UserResponse{
				StatusCode: strings.UnableToQueryUserErrorCode,
				StatusMsg:  strings.UnableToQueryUserError,
			}, nil
		}
		// 其他错误，直接返回错误
		return nil, err
	}

	// 返回用户信息
	resp = &user.UserResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		User:       &queryUser,
	}

	return resp, nil
}

// 根据 username 查询用户是否存在
func CheckUserExists(username string) (*user.UserExistResponse, error) {
	collection := database.MongoDbClient.Database("aorb").Collection("users")
	filter := bson.M{"username": username}
	var queryUser user.User
	err := collection.FindOne(context.TODO(), filter).Decode(&queryUser)
	if err != nil {
		// 当用户不存在的时候，返回 false ，没有错误
		if err == mongo.ErrNoDocuments {
			return &user.UserExistResponse{
				StatusCode: strings.ServiceOKCode,
				StatusMsg:  strings.ServiceOK,
				Existed:    false,
			}, nil
		}
		// 当在查询中出现其他的错误时，返回错误
		return nil, err
	}

	return &user.UserExistResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Existed:    true,
	}, nil
}

// 查询一个用户是否关注另外一个用户
func IsUserFollowing(username string, target_username string) (*user.IsUserFollowingResponse, error) {
	collection := database.MongoDbClient.Database("aorb").Collection("users")
	filter := bson.M{"username": username}
	var queryUser user.User
	err := collection.FindOne(context.TODO(), filter).Decode(&queryUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &user.IsUserFollowingResponse{
				StatusCode:  strings.ServiceOKCode,
				StatusMsg:   strings.ServiceOK,
				IsFollowing: false,
			}, nil
		}
		return nil, err
	}

	// 检查 followed.userids 中是否包含 target_username
	isFollowing := false
	for _, username := range queryUser.Followed.Usernames {
		if username == target_username {
			isFollowing = true
			break
		}
	}

	return &user.IsUserFollowingResponse{
		StatusCode:  strings.ServiceOKCode,
		StatusMsg:   strings.ServiceOK,
		IsFollowing: isFollowing,
	}, nil
}

func UpdateUserInService(ctx context.Context, userId string, updateFields map[string]interface{}) (resp *user.UpdateUserResponse, err error) {
	log.Debug("Updating user in service: ", userId)
	collection := database.MongoDbClient.Database("aorb").Collection("users")

	// 如果是对username进行更新，需要先检查是否有重复的username
	if username, ok := updateFields["username"].(string); ok {
		exists, err := checkUserExistsbyUsername(username)
		if err != nil {
			log.Error("Failed to check username existence: ", err)
			return nil, err
		}

		// 如果存在同名用户，返回错误，到最顶层处理
		if exists {
			log.Warn("Username already exists: ", username)
			return nil, errors.New("username already exists")
		}
	}

	// 如果是对图片（avatar/bgic_me/bgpic_pollcard ）进行更新
	// 需要先在数据库中根据url查找deletion，然后进行删除，然后再更新
	fieldsToCheck := []string{"avatar", "bgpic_me", "bgpic_pollcard"}
	// 这里添加了一个 map 映射，在注册的时候mongodb中存的就是bgpicme，但是在proto中是bgpic_me，json映射也是写的bgpic_me
	mongodb_fields := map[string]string{
		"avatar":         "avatar",
		"bgpic_me":       "bgpicme",
		"bgpic_pollcard": "bgpicpollcard",
	}
	for _, field := range fieldsToCheck {
		if _, ok := updateFields[field]; ok {
			err := deleteImage(userId, field, updateFields[field].(*user.SmmsResponse)) // 类型断言
			if err != nil {
				log.Error("Failed to delete image: ", err)
				return nil, err
			}
			result, err := collection.UpdateOne(ctx, bson.M{"id": userId}, bson.M{"$set": bson.M{mongodb_fields[field]: updateFields[field].(*user.SmmsResponse).Url}})
			if err != nil {
				log.Error("Failed to update user: ", err)
				return nil, err
			}

			// 检查是否有文档被更新
			if result.ModifiedCount == 0 {
				log.Warn("No user was updated")
				return &user.UpdateUserResponse{
					StatusCode: strings.NoUserUpdatedErrorCode,
					StatusMsg:  strings.NoUserUpdatedError,
				}, nil
			}

			return &user.UpdateUserResponse{
				StatusCode: strings.ServiceOKCode,
				StatusMsg:  strings.ServiceOK,
			}, nil
		}
	}

	// 非图片的使用这个更新语句
	// 使用 userId 作为过滤条件, 这里只更新了用户信息
	result, err := collection.UpdateOne(ctx, bson.M{"id": userId}, bson.M{"$set": updateFields})
	if err != nil {
		log.Error("Failed to update user: ", err)
		return nil, err
	}

	// 检查是否有文档被更新
	if result.ModifiedCount == 0 {
		log.Warn("No user was updated")
		return &user.UpdateUserResponse{
			StatusCode: strings.NoUserUpdatedErrorCode,
			StatusMsg:  strings.NoUserUpdatedError,
		}, nil
	}

	return &user.UpdateUserResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
	}, nil
}

// checkUserExistsbyUsername 检查用户名是否存在，内部方法
func checkUserExistsbyUsername(username string) (bool, error) {
	collection := database.MongoDbClient.Database("aorb").Collection("users")
	filter := bson.M{"username": username}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		log.Error("Failed to check username existence: ", err)
		return false, err
	}
	return count > 0, nil
}

// deleteImage 删除 sm.ms 图床上的图片
// 根据 userid 查询对应的文档，然后删除 $field 字段的图片
func deleteImage(userId, field string, picdata *user.SmmsResponse) error {
	// 如果是默认的图片，直接返回
	if picdata.Url == conf.DefaultBgpicMe || picdata.Url == conf.DefaultBgpicPollcard {
		return nil
	}

	// 根据userid和field进行查询旧的图片的信息
	collection := database.MongoDbClient.Database("aorb").Collection("pictures")
	filter := bson.M{"userid": userId, "type": field}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Error("Failed to delete image: ", err)
		return err
	}

	// 获取 delete 字段的值
	deleteLink, ok := result["delete"].(string)
	if !ok {
		log.Fatal("delete field is not a string or does not exist")
	}

	// 在 sm.ms 图床上删除图片
	resp, err := http.Get(deleteLink)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 在pictures数据库中删除图片信息
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Error("Failed to delete image: ", err)
		return err
	}
	if deleteResult.DeletedCount == 0 {
		log.Warn("No image was deleted")
		return errors.New("no image was deleted")
	}
	log.Infof("Delete link response status: %s\n", resp.Status)

	// 从 picdata 中获取新的图片信息并写入到数据库中
	_, err = collection.InsertOne(context.TODO(), bson.M{
		"userid": userId,
		"type":   field,
		"url":    picdata.Url,
		"delete": picdata.Delete,
		"hash":   picdata.Hash})
	if err != nil {
		log.Error("Failed to insert new image: ", err)
		return err
	}

	return nil
}
