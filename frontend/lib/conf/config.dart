// config.dart
import 'package:logger/logger.dart';

// api base url
const String apiDomain = '127.0.0.1';

// 后端网关主机地址和监听端口
const String backendHost = 'localhost';
const int backendPort = 37000;

// logger
Logger getLogger() {
  return Logger(
    printer: PrettyPrinter(),
    filter: ProductionFilter(),
  );
}