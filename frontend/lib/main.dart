import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'routes/router.dart';
import 'package:logger/logger.dart';

import 'package:aorb/conf/config.dart';

void main() {
  var logger = getLogger();
  
  // 设置状态栏透明
  SystemChrome.setSystemUIOverlayStyle(
    SystemUiOverlayStyle(
      statusBarColor: Colors.transparent, // 设置状态栏透明
      statusBarIconBrightness: Brightness.dark, // 设置状态栏图标的亮度
    ),
  );

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
