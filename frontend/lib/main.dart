import 'package:flutter/material.dart';
import 'screens/main_screen.dart';
import 'screens/help_me_choose.dart';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
        title: 'Aorb: Decision Helper',
        theme: ThemeData(
          primarySwatch: Colors.blue,
          visualDensity: VisualDensity.adaptivePlatformDensity,
          appBarTheme: const AppBarTheme(
            backgroundColor: Colors.white, // 设置AppBar的背景颜色
            foregroundColor: Colors.black, // 设置AppBar中前景色（如文字、图标）
            centerTitle: true, // 将标题居中
            elevation: 0, // 移除阴影
          ),
        ),
        home: const MainScreen(),
        routes: {
          '/help_me_choose': (context) => const HelpMeChoose(),
        });
  }
}
