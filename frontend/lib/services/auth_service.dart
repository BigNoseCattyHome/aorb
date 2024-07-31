// auth_service.dart
import 'package:aorb/conf/config.dart';
import 'package:grpc/grpc.dart';
import 'package:aorb/generated/auth.pbgrpc.dart';
import 'package:shared_preferences/shared_preferences.dart';

// @sirius2alpha 2024-07-17
class AuthService {
  late final AuthServiceClient _client;
  late final ClientChannel _channel;
  var logger = getLogger();

  // 初始化_channel和_client
  AuthService() {
    const host = backendHost;
    const port = backendPort;
    logger.i('Attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(
          credentials: ChannelCredentials.insecure()), // ! 生产环境需要更改
    );
    _client = AuthServiceClient(_channel);
  }

  Future<LoginResponse> login(LoginRequest request) async {
    try {
      final LoginResponse response = await _client.login(request);
      logger.i('Login response: $response');

      if (response.statusCode == 0) {
        logger.i(
            'Login successful for user: ${request.username.substring(0, 2)}...'); // 部分隐藏用户名

        // 存储用户信息到本地
        final prefs = await SharedPreferences.getInstance();
        await Future.wait([
          prefs.setString('authToken', response.token),
          prefs.setString('refreshToken', response.refreshToken),
          prefs.setString('username', response.simpleUser.username),
          prefs.setString('avatar', response.simpleUser.avatar),
          prefs.setString('nickname', response.simpleUser.nickname),
        ]);
        logger.d('Token saved to local storage: ${response.token}');
        logger.d('Refresh token saved to local storage: ${response.refreshToken}');
        logger.d('Username saved to local storage: ${response.simpleUser.username}');
        logger.d('Avatar saved to local storage: ${response.simpleUser.avatar}');
        logger.d('Nickname saved to local storage: ${response.simpleUser.nickname}');
        
        return response;
      } else {
        logger.w('Login failed: ${response.statusMsg}');
        throw Exception('Login failed: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('gRPC error during login: ${e.message}');
      throw Exception('gRPC error during login: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error during login: $e');
      throw Exception('Unexpected error during login: $e');
    }
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

  // 检查登录状态
  Future<bool> checkLoginStatus() async {
    // 从本地获取token
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString('authToken');
    logger.d('Token from local storage: $token');

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
