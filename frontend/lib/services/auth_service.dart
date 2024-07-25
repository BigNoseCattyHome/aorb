// auth_service.dart
import 'package:aorb/conf/config.dart';
import 'package:grpc/grpc.dart';
import 'package:aorb/generated/auth.pbgrpc.dart';
import 'package:shared_preferences/shared_preferences.dart';

// @sirius2alpha 2024-07-17
class AuthService {
  late final AuthServiceClient _client;
  late final ClientChannel _channel;

  // 初始化_channel和_client
  AuthService(String host, int port) {
    logger.i('Attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(
          credentials: ChannelCredentials.insecure()), // ! 生产环境需要更改
    );
    _client = AuthServiceClient(_channel);
  }

  Future<LoginResponse> login(
      String username, String password, String deviceId) async {
    final request = LoginRequest()
      ..username = username // 相当于 request.username = username
      ..password = password
      ..deviceId = deviceId;
    final LoginResponse response = await _client.login(request);
    try {
      if (response.statusCode == 0) {
        logger.i('Login successful');

        // 存储用户信息到本地
        final prefs = await SharedPreferences.getInstance();
        await prefs.setString('authToken', response.token);
        await prefs.setString('refreshToken', response.refreshToken);
        await prefs.setString('userId', response.simpleUser.username);
        await prefs.setString('avatar', response.simpleUser.avatar);
        await prefs.setString('nickname', response.simpleUser.nickname);
      } else {
        logger.w('Login failed: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      // 处理gRPC错误
      logger.e('gRPC error during login: ${e.message}');
      throw Exception('Failed to login: ${e.message}');
    } catch (e) {
      // 处理其他错误
      logger.e('Unexpected error during login: $e');
      throw Exception('Failed to login: $e');
    }
    return response;
  }

  Future<VerifyResponse> verify(String token) async {
    final request = VerifyRequest()..token = token;
    return await _client.verify(request);
  }

  Future<RefreshResponse> refresh(String refreshToken) async {
    final request = RefreshRequest()..refreshToken = refreshToken;
    return await _client.refresh(request);
  }

  // 退出登录
  Future<LogoutResponse> logout(String accessToken, String refreshToken) async {
    final request = LogoutRequest()
      ..accessToken = accessToken
      ..refreshToken = refreshToken;
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('authToken');
    return await _client.logout(request);
  }

  Future<RegisterResponse> register(String username, String password,
      {String? nickname, String? avatar, String? ipaddress}) async {
    final request = RegisterRequest()
      ..username = username
      ..password = password;

    if (nickname != null) request.nickname = nickname;
    if (avatar != null) request.avatar = avatar;
    if (ipaddress != null) request.ipaddress = ipaddress;

    return await _client.register(request);
  }

  Future<void> dispose() async {
    await _channel.shutdown();
  }

  var logger = getLogger();

  // 检查登录状态
  Future<bool> checkLoginStatus() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('authToken');

    // 这里需要把成功登录返回的user信息返回到MainPage
    if (token != null) {
      return true;
    } else {
      return false;
    }
  }
}

// void main() async {
//   final authService = AuthService('localhost', 9000);

//   try {
//     // 登录
//     final loginResponse =
//         await authService.login('user_id', 'password', 'device_id');
//     print('Login success: ${loginResponse.success}');
//     print('Token: ${loginResponse.token}');

//     // 验证 token
//     final verifyResponse = await authService.verify(loginResponse.token);
//     print('Verify success: ${verifyResponse.success}');

//     // 刷新 token
//     final refreshResponse =
//         await authService.refresh(loginResponse.refreshToken);
//     print('Refresh success: ${refreshResponse.success}');
//     print('New token: ${refreshResponse.token}');

//     // 注册新用户
//     final registerResponse = await authService
//         .register('new_username', 'new_password', nickname: 'New User');
//     print('Register success: ${registerResponse.success}');

//     // 登出
//     final logoutResponse = await authService.logout(
//         loginResponse.token, loginResponse.refreshToken);
//     print('Logout success: ${logoutResponse.success}');
//   } catch (e) {
//     print('Error: $e');
//   } finally {
//     // 关闭连接
//     await authService.dispose();
//   }
// }
