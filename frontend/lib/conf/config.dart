// config.dart
import 'package:logger/logger.dart';

// api base url
const String apiDomain = '127.0.0.1';

// 后端网关主机地址和监听端口
// ! 使用安卓机调试的时候，需要将手机和电脑置于同一局域网下，并填写电脑的IP地址
const String backendHost = '192.168.124.10'; 
const int backendPort = 37000;

// logger
Logger getLogger() {
  return Logger(
    printer: PrettyPrinter(),
    filter: ProductionFilter(),
  );
}