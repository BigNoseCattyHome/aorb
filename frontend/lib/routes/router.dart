import 'package:aorb/screens/settings_page.dart';
import 'package:flutter/material.dart';
import 'package:aorb/screens/main_page.dart';

class AppRouter {
  static const String homeRoute = '/'; // 首页推荐页面
  static const String messagesRoute = '/messages'; // 消息页面
  static const String meRoute = '/me'; // “我”页面
  static const String settingsRoute = '/settings'; // 设置页面

  static Route<dynamic> generateRoute(RouteSettings settings) {
    switch (settings.name) {
      case homeRoute:
        return MaterialPageRoute(
            builder: (_) => const MainPage(initialIndex: 0));
      case messagesRoute:
        return MaterialPageRoute(
            builder: (_) => const MainPage(initialIndex: 1));
      case meRoute:
        return MaterialPageRoute(
            builder: (_) => const MainPage(initialIndex: 2));
      case settingsRoute:
        return MaterialPageRoute(builder: (_) => SettingsPage());
      default:
        return MaterialPageRoute(
            builder: (_) => Scaffold(
                  body: Center(
                    child: Text('No route defined for ${settings.name}'),
                  ),
                ));
    }
  }
}
