import 'package:aorb/conf/config.dart';
import 'package:aorb/models/poll.dart';
import 'package:http/http.dart' as http;
import 'dart:convert';

class PollService {
  
  // 查询单个投票详情
  Future<Poll> fetchPoll(String pollId) async {
    // 通过 pollId 查询投票详情
    final response =
        await http.get(Uri.http(apiDomain, '/api/v1/poll/$pollId'));

    if (response.statusCode == 200) {
      // 如果服务器返回一个 200 OK 响应，那么解析 JSON。
      final data = jsonDecode(response.body);
      return Poll.fromJson(data);
    } else {
      // 如果服务器返回一个不是 200 OK 的响应，那么抛出一个异常。
      throw Exception('Failed to load poll');
    }
  }
}
