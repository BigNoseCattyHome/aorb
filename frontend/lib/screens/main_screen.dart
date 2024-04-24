import 'package:flutter/material.dart';
import 'help_me_choose.dart';
import 'help_them_choose.dart';
import 'messages.dart';
import 'profile.dart';
import '../widgets/bottom_nav_bar.dart';

class MainScreen extends StatefulWidget {
  @override
  _MainScreenState createState() => _MainScreenState();
}

class _MainScreenState extends State<MainScreen> {
  int _currentIndex = 0;
  final List<Widget> _children = [
    HelpMeChoose(),
    HelpThemChoose(),
    Messages(),
    Profile()
  ];

  void onTabTapped(int index) {
    setState(() {
      _currentIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: _children[_currentIndex],
      bottomNavigationBar: BottomNavBar(onTabTapped: onTabTapped, currentIndex: _currentIndex),
    );
  }
}
