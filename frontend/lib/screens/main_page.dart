import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'home_page.dart';
import 'messages_page.dart';
import 'me_page.dart';
import '../widgets/top_bar_index.dart';

class MainPage extends StatefulWidget {
  final int initialIndex;
  const MainPage({super.key, this.initialIndex = 0});

  @override
  MainPageState createState() => MainPageState();
}

class MainPageState extends State<MainPage>
    with SingleTickerProviderStateMixin {
  late int _currentIndex; // 用于控制底部到行栏的切换
  late TabController tabController; // 顶部导航栏控制器

  late final List<Widget> _pages = [
    HomePage(tabController: tabController), // 传递 _tabController 给 HomePage
    MessagesPage(tabController: tabController),
    const MePage(),
  ];

  @override
  void initState() {
    super.initState();
    // vsync: this 表示使用当前的 SingleTickerProviderStateMixin
    tabController = TabController(length: 2, vsync: this);
    _currentIndex = widget.initialIndex;
  }

  // 底部导航栏切换
  void _onItemTapped(int index) {
    setState(() {
      _currentIndex = index;
      // 第三个页面不需要切换标签
      if (_currentIndex == 2) {
        tabController.index = 0;
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // appBar随底部栏进行切换
      appBar: _currentIndex == 2
          // 第三个页面不需要切换标签
          ? null
          : DynamicTopBar(
              tabs:
                  _currentIndex == 0 ? const ['推荐', '关注'] : const ['提醒', '私信'],
              showSearch: true,
              tabController: tabController,
            ),

      // body的内容由 _pages 控制
      body: IndexedStack(
        index: _currentIndex,
        children: _pages,
      ),

      // 侧栏
      drawer: Drawer(
        child: ListView(
          padding: EdgeInsets.zero,
          children: <Widget>[
            ListTile(
              leading: const Icon(Icons.settings),
              title: const Text('设置'),
              onTap: () {
                // Update the state of the app
                // ...
                Navigator.pop(context);
              },
            ),
            
          ],
        ),
      ),
      
      // 底部导航栏
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
        selectedLabelStyle: const TextStyle(
            fontFamily: 'SimHei', fontSize: 12, fontWeight: FontWeight.w700),
        unselectedLabelStyle: const TextStyle(
            fontFamily: 'SimHei', fontSize: 12, fontWeight: FontWeight.w500),
        onTap: _onItemTapped,
      ),
    );
  }

  @override
  void dispose() {
    tabController.dispose();
    super.dispose();
  }
}
