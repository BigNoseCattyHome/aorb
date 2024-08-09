import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:intl/intl.dart';

// 格式化时间戳，第一个参数是时间戳，第二个参数是前缀（可选）
String formatTimestamp(Timestamp timestamp, [String? prefix]) {
  final now = DateTime.now();
  final dateTime = timestamp.toDateTime();
  final difference = now.difference(dateTime);

  if (difference.inMinutes < 60) {
    return '$prefix${difference.inMinutes}分钟前';
  } else if (difference.inHours < 24) {
    return '$prefix${difference.inHours}小时前';
  } else if (difference.inDays < 2) {
    return '$prefix昨天${DateFormat('HH:mm').format(dateTime)}';
  } else if (difference.inDays < 7) {
    return '$prefix${difference.inDays}天前';
  } else if (difference.inDays < 30) {
    return '$prefix${(difference.inDays / 7).floor()}周前';
  } else if (difference.inDays < 365) {
    return '$prefix${(difference.inDays / 30).floor()}个月前';
  } else {
    return '$prefix${DateFormat('yyyy年MM月dd日').format(dateTime)}';
  }
}
