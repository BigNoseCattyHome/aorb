import 'package:aorb/conf/config.dart';
import 'package:aorb/utils/constant/err.dart';
import 'package:grpc/grpc.dart';
import 'package:aorb/generated/user.pbgrpc.dart';

class UserService {
  late final UserServiceClient _client;
  late final ClientChannel _channel;
  var logger = getLogger();

  UserService() {
    const host = backendHost;
    const port = backendPort;
    logger.i('Attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    _client = UserServiceClient(_channel);
  }

  // 查询用户信息
  Future<UserResponse> getUserInfo(UserRequest request) async {
    try {
      final UserResponse response = await _client.getUserInfo(request);

      if (response.statusCode == 0) {
        logger.i(
            'Successfully retrieved information for user ${request.username}');
        return response;
      } else {
        logger.w(
            'Failed to retrieve information for user ${request.username}: ${response.statusMsg}');
        throw Exception(
            'Failed to retrieve information for user ${request.username}: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e(
          'gRPC error occurred while retrieving user information: ${e.message}');
      throw Exception('Failed to retrieve user information: ${e.message}');
    } catch (e) {
      logger
          .e('Unexpected error occurred while retrieving user information: $e');
      throw Exception('Failed to retrieve user information: $e');
    }
  }

  // 查询用户是否存在
  Future<UserExistResponse> checkUserExists(UserExistRequest request) async {
    try {
      final UserExistResponse response = await _client.checkUserExists(request);

      if (response.statusCode == 0) {
        logger.i('Successfully checked if user ${request.username} exists');
        return response;
      } else {
        logger.w(
            'Failed to check if user ${request.username} exists: ${response.statusMsg}');
        throw Exception(
            'Failed to check if user ${request.username} exists: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger
          .e('gRPC error occurred while checking if user exists: ${e.message}');
      throw Exception('Failed to check if user exists: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error occurred while checking if user exists: $e');
      throw Exception('Failed to check if user exists: $e');
    }
  }

  // 查询一个用户是否关注另一个用户
  Future<bool> isUserFollowing(IsUserFollowingRequest request) async {
    try {
      final IsUserFollowingResponse response =
          await _client.isUserFollowing(request);

      if (response.statusCode == 0) {
        logger.i(
            'Successfully checked if user ${request.username} is following ${request.targetUsername}');
        return response.isFollowing;
      } else {
        logger.w(
            'Failed to check if user ${request.username} is following ${request.targetUsername}: ${response.statusMsg}');
        throw Exception(
            'Failed to check if user ${request.username} is following ${request.targetUsername}: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e(
          'gRPC error occurred while checking user follow status: ${e.message}');
      throw Exception('Failed to check user follow status: ${e.message}');
    } catch (e) {
      logger
          .e('Unexpected error occurred while checking user follow status: $e');
      throw Exception('Failed to check user follow status: $e');
    }
  }

  // 编辑用户资料实现用户信息更新
  // 图片在调用方法的时候上传图床，只用传递一个 UpdateUserRequest 对象
  Future<UpdateUserResponse> updateUser(UpdateUserRequest request) async {
    try {
      final UpdateUserResponse response = await _client.updateUser(request);
      if (response.statusCode == 0) {
        logger.i('Successfully updated information for user ${request.userId}');
      } else if (response.statusCode == authUserExistedCode) {
        logger.w('Username ${request.username} already exists');
      } else {
        logger.w(
            'Failed to update information for user ${request.username}: ${response.statusMsg}');
      }
      return response;
    } on GrpcError catch (e) {
      logger.e(
          'gRPC error occurred while updating user information: ${e.message}');
      return UpdateUserResponse()
        ..statusCode = -1
        ..statusMsg = 'Failed to update user information: ${e.message}';
    } catch (e) {
      logger.e('Unexpected error occurred while updating user information: $e');
      return UpdateUserResponse()
        ..statusCode = -1
        ..statusMsg = 'Failed to update user information: $e';
    }
  }

  // 关注用户
  Future<FollowUserResponse> followUser(FollowUserRequest request) async {
    try {
      final FollowUserResponse response = await _client.followUser(request);

      if (response.statusCode == 0) {
        logger.i(
            'Successfully followed user ${request.targetUsername} as user ${request.username}');
      } else {
        logger.w(
            'Failed to follow user ${request.targetUsername} as user ${request.username}: ${response.statusMsg}');
        throw Exception(
            'Failed to follow user ${request.targetUsername} as user ${request.username}: ${response.statusMsg}');
      }
      return response;
    } on GrpcError catch (e) {
      logger.e('gRPC error occurred while following user: ${e.message}');
      throw Exception('Failed to follow user: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error occurred while following user: $e');
      throw Exception('Failed to follow user: $e');
    }
  }

  // 取消关注用户
  Future<FollowUserResponse> unfollowUser(FollowUserRequest request) async {
    try {
      final FollowUserResponse response = await _client.unfollowUser(request);

      if (response.statusCode == 0) {
        logger.i(
            'Successfully unfollowed user ${request.targetUsername} as user ${request.username}');
      } else {
        logger.w(
            'Failed to unfollow user ${request.targetUsername} as user ${request.username}: ${response.statusMsg}');
        throw Exception(
            'Failed to unfollow user ${request.targetUsername} as user ${request.username}: ${response.statusMsg}');
      }
      return response;
    } on GrpcError catch (e) {
      logger.e('gRPC error occurred while unfollowing user: ${e.message}');
      throw Exception('Failed to unfollow user: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error occurred while unfollowing user: $e');
      throw Exception('Failed to unfollow user: $e');
    }
  }
}
