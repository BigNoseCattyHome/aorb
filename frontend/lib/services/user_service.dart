import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:aorb/conf/config.dart';

class UserService {
  // 查询一个用户是否关注另一个用户
  Future<bool> fetchFollowStatus(String userId, String someoneUserId) async {
    final response =
        await http.get(Uri.http(apiDomain, '/api/v1/user/$userId', {
      'someone_user_id': someoneUserId,
    }));

    if (response.statusCode == 200) {
      // 如果服务器返回一个 200 OK 响应，那么解析 JSON。
      final data = jsonDecode(response.body);
      return data['isFollowed'];
    } else {
      // 如果服务器返回一个不是 200 OK 的响应，那么抛出一个异常。
      throw Exception('Failed to load follow status');
    }
  }

  // 查询用户的部分或者全部信息
  Future<Map<String, dynamic>> fetchUserInfo(String userId,
      [List<String> fields = const []]) async {
    // 构建查询参数
    String queryParams = fields.isNotEmpty ? '?fields=${fields.join(',')}' : '';

    // 发送GET请求
    final response =
        await http.get(Uri.http(apiDomain, '/api/v1/user/$userId$queryParams'));

    if (response.statusCode == 200) {
      // 如果服务器返回一个 200 OK 响应，那么解析 JSON。
      final data = jsonDecode(response.body);
      return data; // ! 注意这个data的解析过程，并不是一个完整的user对象
    } else {
      // 如果服务器返回一个不是 200 OK 的响应，那么抛出一个异常。
      throw Exception('Failed to load user info');
    }
  }
}
