import 'package:flutter/material.dart';
import 'routes/router.dart';
import 'package:logger/logger.dart';

import 'package:aorb/conf/config.dart';

void main() {
  var logger = getLogger();
  runApp(Aorb(logger: logger));
}

class Aorb extends StatelessWidget {
  final Logger logger;

  const Aorb({super.key, required this.logger});

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
