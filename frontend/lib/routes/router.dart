import 'package:flutter/material.dart';
import 'package:aorb/screens/main_page.dart';
import 'package:aorb/screens/poll_detail_page.dart';

class AppRouter {
  static const String homeRoute = '/'; // 首页推荐页面
  static const String messagesRoute = '/messages'; // 消息页面
  static const String meRoute = '/me'; // “我”页面
  static const String pollContentRoute = '/poll_content'; // 内容详情页面

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
      case pollContentRoute:
        // 假设你从 RouteSettings 中获取了 postUserID 和 userID
        final postUserId = settings.arguments as String;
        final userId = settings.arguments as String;
        return MaterialPageRoute(
            builder: (_) =>
                PollDetailPage(postUserId: postUserId, username: userId));
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
