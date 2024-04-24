import 'package:flutter/material.dart';

class BottomNavBar extends StatelessWidget {
  final int currentIndex;
  final Function(int) onTabTapped;

  BottomNavBar({required this.onTabTapped, required this.currentIndex});

  @override
  Widget build(BuildContext context) {
    return BottomNavigationBar(
      backgroundColor: Colors.blue, // 设置底部导航栏的背景颜色为蓝色
      selectedItemColor: Colors.blue, // 设置选中项的颜色为白色
      unselectedItemColor: Colors.grey, // 设置未选中项的颜色为灰色
      onTap: onTabTapped,
      currentIndex: currentIndex,
      items: const [
        BottomNavigationBarItem(icon: Icon(Icons.help_outline), label: '帮我选'),
        BottomNavigationBarItem(icon: Icon(Icons.people_outline), label: '帮他选'),
        BottomNavigationBarItem(icon: Icon(Icons.message), label: '消息'),
        BottomNavigationBarItem(icon: Icon(Icons.person), label: '我的'),
      ],
    );
  }
}
