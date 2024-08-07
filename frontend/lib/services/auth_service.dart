// auth_service.dart
import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/user.pb.dart';
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
        credentials: ChannelCredentials.insecure(), // ! 生产环境需要更改
        idleTimeout: Duration(seconds: 60), // 增加空闲超时
      ),
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
        logger.d(
            'Refresh token saved to local storage: ${response.refreshToken}');
        logger.d(
            'Username saved to local storage: ${response.simpleUser.username}');
        logger
            .d('Avatar saved to local storage: ${response.simpleUser.avatar}');
        logger.d(
            'Nickname saved to local storage: ${response.simpleUser.nickname}');

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

  Future<RegisterResponse> register(
      String username, String password, Gender gender,
      {String? nickname, String? avatar, String? ipaddress}) async {
    final request = RegisterRequest()
      ..username = username
      ..password = password
      ..gender = gender;

    if (nickname != null) request.nickname = nickname;
    if (avatar != null) request.avatar = avatar;
    if (ipaddress != null) request.ipaddress = ipaddress;

    return await _client.register(request,
        options: CallOptions(timeout: const Duration(seconds: 30)));
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
