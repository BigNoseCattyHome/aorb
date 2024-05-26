import 'package:flutter/material.dart';
import 'routes/router.dart';

void main() {
  runApp(const Aorb());
}

class Aorb extends StatelessWidget {
  const Aorb({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Aorb: Decision Helper',
      theme: ThemeData(
        scaffoldBackgroundColor: Colors.white,
      ),
      onGenerateRoute: AppRouter.generateRoute, // 使用自定义路由
      initialRoute: AppRouter.homeRoute, // 设置初始路由
    );
  }
}
