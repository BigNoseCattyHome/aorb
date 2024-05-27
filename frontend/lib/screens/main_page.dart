import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'home_page.dart';
import 'messages_page.dart';
import 'me_page.dart';

class MainPage extends StatefulWidget {
  final int initialIndex;
  const MainPage({super.key, this.initialIndex = 0});

  @override
  MainPageState createState() => MainPageState();
}

class MainPageState extends State<MainPage> {
  late int _currentIndex;
  final List<Widget> _pages = [
    // 加了const之后，避免了每次都重新创建对象
    const HomePage(),
    const MessagesPage(),
    const MePage(),
  ];

  @override
  void initState() {
    super.initState();
    _currentIndex = widget.initialIndex;
  }

  void _onItemTapped(int index) {
    setState(() {
      _currentIndex = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(
        index: _currentIndex,
        children: _pages,
      ),
      bottomNavigationBar: BottomNavigationBar(
        items: <BottomNavigationBarItem>[
          BottomNavigationBarItem(
            icon: _currentIndex == 0
                ? SvgPicture.asset('images/home_selected.svg')
                : SvgPicture.asset('images/home_unselected.svg'),
            label: '首页',
          ),
          BottomNavigationBarItem(
            icon: _currentIndex == 1
                ? SvgPicture.asset('images/msg_selected.svg')
                : SvgPicture.asset('images/msg_unselected.svg'),
            label: '消息',
          ),
          BottomNavigationBarItem(
            icon: _currentIndex == 2
                ? SvgPicture.asset('images/me_selected.svg')
                : SvgPicture.asset('images/me_unselected.svg'),
            label: '我',
          ),
        ],
        currentIndex: _currentIndex,
        selectedItemColor: Colors.blue[700],
        selectedLabelStyle: TextStyle(
            fontFamily: 'SimHei', fontSize: 12, fontWeight: FontWeight.w700),
        unselectedLabelStyle: TextStyle(fontFamily: 'SimHei', fontSize: 12,fontWeight: FontWeight.w500),
        onTap: _onItemTapped,
      ),
    );
  }
}
