import 'package:flutter/material.dart';
import 'index.dart';
import 'messages.dart';
import 'profile.dart';
import '../widgets/bottom_nav_bar.dart';

// MainScreen is the main screen of the app, which contains three tabs: Index, Messages, and Profile.
class MainScreen extends StatefulWidget {
  // MainScreen({Key? key}) : super(key: key);
  // 表示MainScreen类的构造函数，接受一个可选的Key类型的参数key
  // super表示调用父类的构造函数，这里调用StatefulWidget类的构造函数，传入key参数
  const MainScreen({super.key});

  // MainScreen重载了createState()方法，返回_MainScreenState的实例
  @override
  _MainScreenState createState() => _MainScreenState();
}

// _MainScreenState类继承自State<MainScreen>类
class _MainScreenState extends State<MainScreen> {
  int _currentIndex = 0;  // 当前选中的tab的索引
  final List<Widget> _children = [  // 三个tab对应的页面
    const Index(),
    Messages(),
    Profile()
  ];

  // 当点击底部导航栏的某个tab时，更新当前选中的tab的索引
  void onTabTapped(int index) {
    setState(() {
      _currentIndex = index;
    });
  }

  // build方法返回一个Scaffold组件，包含一个body和一个bottomNavigationBar
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _children[_currentIndex],
      bottomNavigationBar: BottomNavBar(onTabTapped: onTabTapped, currentIndex: _currentIndex),
    );
  }
}
