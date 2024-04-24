// 目前还没有应用

import 'package:flutter/material.dart';

class BottomNavBarColorController with ChangeNotifier {
  Color _color = Colors.blue; // 默认颜色

  Color get color => _color;

  void setColor(Color newColor) {
    _color = newColor;
    notifyListeners(); // 通知监听器颜色已更改
  }
}
