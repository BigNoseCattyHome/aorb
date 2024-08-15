import 'package:aorb/conf/config.dart';
import 'package:flutter/material.dart';
import 'package:palette_generator/palette_generator.dart';
import 'package:flutter_cache_manager/flutter_cache_manager.dart';

class ColorAnalyzer {
  static final Map<String, Color> _cache = {};
  static final logger = getLogger();

  static Future<Color> getTextColor(String backgroundImage) async {
    if (_cache.containsKey(backgroundImage)) {
      return _cache[backgroundImage]!;
    }

    Color textColor;

    if (backgroundImage.startsWith('0x')) {
      // 纯色背景
      int colorValue = int.parse(backgroundImage.substring(2), radix: 16);
      Color backgroundColor = Color(colorValue);
      textColor = _isLightColor(backgroundColor) ? Colors.black : Colors.white;
    } else if (backgroundImage.startsWith('gradient:')) {
      // 渐变背景，取第一个颜色作为参考
      List<String> colorStrings =
          backgroundImage.substring('gradient:'.length).split(',');
      Color firstColor =
          Color(int.parse(colorStrings[0].substring(2), radix: 16));
      textColor = _isLightColor(firstColor) ? Colors.black : Colors.white;
    } else {
      // 图片背景
      textColor = await _analyzeImageColor(backgroundImage);
    }

    _cache[backgroundImage] = textColor;
    return textColor;
  }

  static Future<Color> _analyzeImageColor(String imageUrl) async {
    try {
      // 使用 CacheManager 来缓存图片
      final file = await DefaultCacheManager().getSingleFile(imageUrl);

      final PaletteGenerator paletteGenerator =
          await PaletteGenerator.fromImageProvider(
        FileImage(file),
        size: const Size(100, 100), // 可以调整大小以提高性能
        maximumColorCount: 20, // 增加色彩分析的精度
      );

      // 尝试获取不同的调色板颜色
      final Color? color = paletteGenerator.dominantColor?.color ??
          paletteGenerator.vibrantColor?.color ??
          paletteGenerator.mutedColor?.color ??
          paletteGenerator.lightVibrantColor?.color ??
          paletteGenerator.darkVibrantColor?.color;

      if (color != null) {
        return _isLightColor(color) ? Colors.black : Colors.white;
      }
    } catch (e) {
      logger.e('Failed to analyze image color: $e');
    }

    // 如果无法分析图片颜色，默认使用黑色文字
    return Colors.black;
  }

  static bool _isLightColor(Color color) {
    return color.computeLuminance() > 0.5;
  }
}
