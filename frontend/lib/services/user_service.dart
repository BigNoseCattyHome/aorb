import 'package:aorb/conf/config.dart';
import 'package:grpc/grpc.dart';
import 'package:aorb/generated/user.pbgrpc.dart';

class UserService {
  late final UserServiceClient _client;
  late final ClientChannel _channel;
  var logger = getLogger();

  // 初始化_channel和_client
  UserService(String host, int port) {
    logger.i('Attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(
          credentials: ChannelCredentials.insecure()), // ! 生产环境需要更改
    );
    _client = UserServiceClient(_channel);
  }

  // 查询用户信息
  Future<UserResponse> getUserInfo(UserRequest request) async {
    try {
      final UserResponse response = await _client.getUserInfo(request);

      // TODO 考虑使用枚举或常量来表示状态码
      if (response.statusCode == 0) {
        logger.i('Get user ${request.username} info successfully');
        return response;
      } else {
        logger.w(
            'Get user ${request.username} info failed: ${response.statusMsg}');
        throw Exception(
            'Get user ${request.username} info failed: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('gRPC error during getting user info: ${e.message}');
      throw Exception('Failed to get user info: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error during getting user info: $e');
      throw Exception('Failed to get user info: $e');
    }
  }

  // 查询用户是否存在
  Future<UserExistResponse> checkUserExists(UserExistRequest request) async {
    try {
      final UserExistResponse response = await _client.checkUserExists(request);

      if (response.statusCode == 0) {
        logger.i('check user ${request.username} exists successfully');
        return response;
      } else {
        logger.w(
            'check user ${request.username} exists failed: ${response.statusMsg}');
        throw Exception(
            'check user ${request.username} exists failed: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('gRPC error during login: ${e.message}');
      throw Exception('Failed to login: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error during login: $e');
      throw Exception('Failed to login: $e');
    }
  }

  // 查询一个用户是否关注另一个用户
  Future<bool> isUserFollowing(IsUserFollowingRequest request) async {
    try {
      final IsUserFollowingResponse response =
          await _client.isUserFollowing(request);

      if (response.statusCode == 0) {
        logger.i(
            'check user ${request.username} is following ${request.targetUsername} successfully');
        return response.isFollowing;
      } else {
        logger.w(
            'check user ${request.username} is following ${request.targetUsername} failed: ${response.statusMsg}');
        throw Exception(
            'check user ${request.username} is following ${request.targetUsername} failed: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('gRPC error during login: ${e.message}');
      throw Exception('Failed to login: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error during login: $e');
      throw Exception('Failed to login: $e');
    }
  }

  // 查询用户的关注列表
  Future<List<User>> getFollowList(String username) async {
    // 创建UserRequest对象
    UserRequest request = UserRequest()
      ..username = username
      ..fields.addAll(['followed']); // 使用addAll方法添加字段

    try {
      // 调用getUserInfo方法，获取用户信息
      UserResponse response = await getUserInfo(request);

      // 判断user对象是否有followed字段
      if (response.user.hasFollowed()) {
        return response.user.followed.userIds
            .map((userId) => User()..id = userId)
            .toList();
      } else {
        return []; // 如果没有关注列表，返回空列表
      }
    } catch (e) {
      logger.e('Error getting follow list: $e');
      throw Exception('Failed to get follow list: $e');
    }
  }

  // 查询用户的粉丝列表
  Future<List<User>> getFanList(String username) async {
    // 创建UserRequest对象
    UserRequest request = UserRequest()
      ..username = username
      ..fields.addAll(['follower']); // 使用addAll方法添加字段

    try {
      // 调用getUserInfo方法，获取用户信息
      UserResponse response = await getUserInfo(request);

      // 判断user对象是否有follower字段
      if (response.user.hasFollower()) {
        return response.user.followed.userIds
            .map((userId) => User()..id = userId)
            .toList();
      } else {
        return []; // 如果没有粉丝列表，返回空列表
      }
    } catch (e) {
      logger.e('Error getting fan list: $e');
      throw Exception('Failed to get fan list: $e');
    }
  }
}
