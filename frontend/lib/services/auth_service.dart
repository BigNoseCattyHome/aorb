// auth_service.dart
import 'package:shared_preferences/shared_preferences.dart';
import 'package:aorb/conf/config.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class AuthService extends Object {
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

  // 登录
  // 这里的password已经是md5摘要了
  Future<void> login(String username, String password) async {
    final response = await http.post(
      Uri.parse('$apiDomain/api/v1/auth/login'),
      headers: <String, String>{
        'Content-Type': 'application/json; charset=UTF-8',
      },
      body: jsonEncode(<String, String>{
        'username': username,
        'password': password,
      }),
    );

    if (response.statusCode == 200) {
      logger.i('Login successful');
      // 把JWT,userid,avtar,nickname存储在本地
      final prefs = await SharedPreferences.getInstance();
      final data = jsonDecode(response.body);
      await prefs.setString('authToken', data['token']);
      await prefs.setString('userId', data['user']['id']);
      await prefs.setString('avatar', data['user']['avatar']);
      await prefs.setString('nickname', data['user']['nickname']);
    } else {
      logger.e('Failed to login');
      throw Exception('Failed to login');
    }
  }

  // 退出登录
  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('authToken');
  }
}
