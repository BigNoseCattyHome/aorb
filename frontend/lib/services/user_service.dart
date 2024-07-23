import 'package:aorb/models/simple_user.dart';
import 'package:aorb/models/user.dart';
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
      final data = jsonDecode(response.body);
      return data['isFollowed'];
    } else {
      throw Exception('Failed to load follow status');
    }
  }

  // 查询用户的部分或者全部信息
  Future<User> fetchUserInfo(String userId,
      [List<String> fields = const []]) async {
    String queryParams = fields.isNotEmpty ? '?fields=${fields.join(',')}' : '';
    final response =
        await http.get(Uri.http(apiDomain, '/api/v1/user/$userId$queryParams'));

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      User user = User.fromJson(data); // 这里把不是很完全的用户信息转换为User对象
      return user;
    } else {
      throw Exception('Failed to load user info');
    }
  }

  // 查询用户的关注列表
  Future<List<SimpleUser>> fetchFollowList(String userId) async {
    final response =
        await http.get(Uri.http(apiDomain, '/api/v1/user/$userId/followed'));

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data;
    } else {
      throw Exception('Failed to load follow list');
    }
  }

  // 查询用户的粉丝列表
  Future<List<SimpleUser>> fetchFanList(String userId) async {
    final response =
        await http.get(Uri.http(apiDomain, '/api/v1/user/$userId/fans'));

    if (response.statusCode == 200) {
      final data = jsonDecode(response.body);
      return data;
    } else {
      throw Exception('Failed to load follow list');
    }
  }
}
