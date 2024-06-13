import 'package:flutter/material.dart';
import 'package:aorb/screens/main_page.dart';
import 'package:aorb/screens/poll_content.dart';

class AppRouter {
  static const String homeRoute = '/';
  static const String messagesRoute = '/messages';
  static const String meRoute = '/me';
  static const String pollContentRoute = '/poll_content'; // 新增的路由

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
        final postUserID = settings.arguments as String;
        final userID = settings.arguments as String;
        return MaterialPageRoute(
            builder: (_) =>
                PollContent(postUserID: postUserID, userID: userID));
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
