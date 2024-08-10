import 'dart:convert';
import 'dart:io';
import 'package:aorb/generated/user.pb.dart';
import 'package:http/http.dart' as http;
import 'package:flutter_dotenv/flutter_dotenv.dart';

class ImageUploadService {
  // SM.MS API URL
  static const String _smmsApiUrl = 'https://sm.ms/api/v2/upload';

  // 上传图片
  Future<smmsResponse> uploadImage(
      File imageFile, String field, String userid) async {
    // 从 .env 文件中获取 SM.MS Token
    final smmsToken = dotenv.env['SMMS_TOKEN'];

    if (smmsToken == null) {
      throw Exception('SMMS_TOKEN not found in .env file');
    }

    // 重命名文件
    final tempDir = await Directory.systemTemp.createTemp();
    final renamedFile =
        await imageFile.copy('${tempDir.path}/${field}_$userid}');

    var request = http.MultipartRequest('POST', Uri.parse(_smmsApiUrl));
    request.headers['Authorization'] = smmsToken;

    request.files
        .add(await http.MultipartFile.fromPath('smfile', renamedFile.path));

    var response = await request.send();
    var responseBody = await response.stream.bytesToString();

    if (response.statusCode == 200) {
      var jsonResponse = json.decode(responseBody);
      if (jsonResponse['success']) {
        var data = jsonResponse['data'];
        return smmsResponse()
          ..url = data['url']
          ..delete = data['delete']
          ..hash = data['hash'];
      } else {
        throw Exception('上传图片失败: ${jsonResponse['message']}');
      }
    } else {
      throw Exception('上传图片失败: ${response.statusCode}');
    }
  }
}
