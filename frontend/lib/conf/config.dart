// config.dart
import 'package:logger/logger.dart';

// api base url
const String apiDomain = '127.0.0.1';

// logger
Logger getLogger() {
  return Logger(
    printer: PrettyPrinter(),
    filter: ProductionFilter(),
  );
}