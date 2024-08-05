import 'package:aorb/conf/config.dart';
import 'package:dart_ipify/dart_ipify.dart';

import 'package:http/http.dart' as http;
import 'dart:convert';

class IPLocationUtil {
  static Future<String> getProvince() async {
    final logger = getLogger();
    try {
      // 获取 IP 地址
      final ip = await Ipify.ipv4();

      // 使用 IP-API 查询地理位置信息
      final response = await http.get(Uri.parse('http://ip-api.com/json/$ip'));
      if (response.statusCode == 200) {
        final data = json.decode(response.body);
        // 返回省份信息
        return data['regionName'] ?? 'Unknown';
      } else {
        return 'Failed to get location';
      }
    } catch (e) {
      logger.e('Error: $e');
      return 'Error occurred';
    }
  }
}
