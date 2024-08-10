import 'package:aorb/conf/config.dart';
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
        logger.i('成功获取用户 ${request.username} 的信息');
        return response;
      } else {
        logger.w('获取用户 ${request.username} 信息失败: ${response.statusMsg}');
        throw Exception('获取用户 ${request.username} 信息失败: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('获取用户信息时发生gRPC错误: ${e.message}');
      throw Exception('获取用户信息失败: ${e.message}');
    } catch (e) {
      logger.e('获取用户信息时发生意外错误: $e');
      throw Exception('获取用户信息失败: $e');
    }
  }

  // 查询用户是否存在
  Future<UserExistResponse> checkUserExists(UserExistRequest request) async {
    try {
      final UserExistResponse response = await _client.checkUserExists(request);

      if (response.statusCode == 0) {
        logger.i('成功检查用户 ${request.username} 是否存在');
        return response;
      } else {
        logger.w('检查用户 ${request.username} 是否存在失败: ${response.statusMsg}');
        throw Exception(
            '检查用户 ${request.username} 是否存在失败: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('检查用户是否存在时发生gRPC错误: ${e.message}');
      throw Exception('检查用户是否存在失败: ${e.message}');
    } catch (e) {
      logger.e('检查用户是否存在时发生意外错误: $e');
      throw Exception('检查用户是否存在失败: $e');
    }
  }

  // 查询一个用户是否关注另一个用户
  Future<bool> isUserFollowing(IsUserFollowingRequest request) async {
    try {
      final IsUserFollowingResponse response =
          await _client.isUserFollowing(request);

      if (response.statusCode == 0) {
        logger.i('成功检查用户 ${request.username} 是否关注 ${request.targetUsername}');
        return response.isFollowing;
      } else {
        logger.w(
            '检查用户 ${request.username} 是否关注 ${request.targetUsername} 失败: ${response.statusMsg}');
        throw Exception(
            '检查用户 ${request.username} 是否关注 ${request.targetUsername} 失败: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('检查用户关注状态时发生gRPC错误: ${e.message}');
      throw Exception('检查用户关注状态失败: ${e.message}');
    } catch (e) {
      logger.e('检查用户关注状态时发生意外错误: $e');
      throw Exception('检查用户关注状态失败: $e');
    }
  }

  // 编辑用户资料实现用户信息更新
  // 图片在调用方法的时候上传图床，只用传递一个 UpdateUserRequest 对象
  Future<UpdateUserResponse> updateUser(UpdateUserRequest request) async {
    try {
      // 调用gRPC接口更新用户信息
      final UpdateUserResponse response = await _client.updateUser(request);

      // 根据返回的状态码判断是否成功
      if (response.statusCode == 0) {
        logger.i('成功编辑用户 ${request.userId} 的信息');
        return response;
      } else {
        logger.w('编辑用户 ${request.username} 信息失败: ${response.statusMsg}');
        throw Exception('编辑用户 ${request.username} 信息失败: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('更新用户信息时发生gRPC错误: ${e.message}');
      throw Exception('编辑用户信息失败: ${e.message}');
    } catch (e) {
      logger.e('更新用户信息时发生意外错误: $e');
      throw Exception('编辑用户信息失败: $e');
    }
  }
}
