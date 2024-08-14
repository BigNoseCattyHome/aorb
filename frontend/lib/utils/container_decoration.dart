import 'package:flutter/material.dart';

// 方法：根据backgroundImage创建背景装饰
BoxDecoration createBackgroundDecoration(String backgroundImage) {
  if (backgroundImage.startsWith('0x')) {
    // 纯色背景
    int colorValue = int.parse(backgroundImage.substring(2), radix: 16);
    return BoxDecoration(
      color: Color(colorValue),
      borderRadius: BorderRadius.circular(10),
    );
  } else if (backgroundImage.startsWith('http://') ||
      backgroundImage.startsWith('https://')) {
    // 网络图片背景
    return BoxDecoration(
      image: DecorationImage(
        image: NetworkImage(backgroundImage),
        fit: BoxFit.cover,
      ),
      borderRadius: BorderRadius.circular(10),
    );
  } else if (backgroundImage.startsWith('gradient:')) {
    // 渐变背景
    List<String> colorStrings =
        backgroundImage.substring('gradient:'.length).split(',');
    List<Color> colors = colorStrings
        .map((colorString) =>
            Color(int.parse(colorString.substring(2), radix: 16)))
        .toList();
    return BoxDecoration(
      gradient: LinearGradient(
        colors: colors,
        begin: Alignment.topLeft,
        end: Alignment.bottomRight,
      ),
      borderRadius: BorderRadius.circular(10),
    );
  } else {
    // 默认纯色背景
    return BoxDecoration(
      color: Colors.blue[200],
      borderRadius: BorderRadius.circular(10),
    );
  }
}
