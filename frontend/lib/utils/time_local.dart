import 'package:timezone/timezone.dart' as tz;
import 'package:timezone/data/latest.dart' as tz;
import 'package:timezone/timezone.dart';

TZDateTime toLocalTime(DateTime time) {
  // 初始化时区数据
  tz.initializeTimeZones();
  // 设置东八区
  tz.setLocalLocation(tz.getLocation('Asia/Shanghai'));

  // 转换时间为东八区时间
  var localTime = tz.TZDateTime.from(time, tz.local);
  return localTime;
}
