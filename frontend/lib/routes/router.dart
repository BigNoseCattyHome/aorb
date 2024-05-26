import 'package:flutter/material.dart';
import '../screens/main_page.dart'; 

class AppRouter {
  static const String homeRoute = '/';
  static const String messagesRoute = '/messages';
  static const String meRoute = '/me';

  static Route<dynamic> generateRoute(RouteSettings settings) {
    switch (settings.name) {
      case homeRoute:
        return MaterialPageRoute(builder: (_) => const MainPage(initialIndex: 0));
      case messagesRoute:
        return MaterialPageRoute(builder: (_) => const MainPage(initialIndex: 1));
      case meRoute:
        return MaterialPageRoute(builder: (_) => const MainPage(initialIndex: 2));
      default:
        return MaterialPageRoute(builder: (_) => Scaffold(
              body: Center(
                child: Text('No route defined for ${settings.name}'),
              ),
            ));
    }
  }
}
