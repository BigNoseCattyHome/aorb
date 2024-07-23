import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';

import 'package:aorb/screens/home_page.dart';
import 'package:aorb/screens/messages_page.dart';
import 'package:aorb/screens/me_page.dart';
import 'package:aorb/screens/login_page.dart';
import 'package:aorb/widgets/top_bar_index.dart';
import 'package:aorb/services/auth_service.dart';
import 'package:shared_preferences/shared_preferences.dart';

class MainPage extends StatefulWidget {
  final int initialIndex;
  const MainPage({super.key, this.initialIndex = 0});
  @override
  MainPageState createState() => MainPageState();
}

class MainPageState extends State<MainPage>
    with SingleTickerProviderStateMixin {
  late int _currentIndex; // 用于控制底部到行栏的切换
  late TabController tabController; // tabController用于控制子页面的顶部导航栏的切换
  late bool isLoggedIn; // 是否登录
  late String userId; // 在initstate中从本地读取，用于传递给 _pages中的MePage
  late String avatar; // 在initstate中从本地读取，用于底部状态栏的icon的展示
  late List<Widget> _pages;

  // 将异步初始化逻辑移到单独的方法中
  Future<void> _initializeData() async {
    try {
      bool loginStatus = await AuthService('localhost', 9000).checkLoginStatus();
      setState(() {
        isLoggedIn = loginStatus;
      });

      if (isLoggedIn) {
        SharedPreferences prefs = await SharedPreferences.getInstance();
        setState(() {
          userId = prefs.getString('userId') ?? '';
          avatar = prefs.getString('avatar') ?? '';
          _pages[2] = MePage(userId: userId);
        });
      }
    } catch (e) {
      print("Error on mainpage: $e");
      // 可以在这里处理错误,例如显示一个错误提示
    }
  }

  @override
  void initState() {
    super.initState();
    // vsync是一个同步的信号，用于保证动画的同步性
    tabController = TabController(length: 2, vsync: this);
    _currentIndex = widget.initialIndex;

    // 提供初始值
    setState(() {
      userId = '';
      avatar = '';
      _pages = [
        HomePage(tabController: tabController),
        MessagesPage(tabController: tabController),
        LoginPage(),
      ];
    });

    // 异步获取登录状态
    _initializeData();
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
                ? SvgPicture.asset(
                    'images/home_selected.svg',
                    width: 35,
                    height: 35,
                    fit: BoxFit.contain,
                  )
                : SvgPicture.asset(
                    'images/home_unselected.svg',
                    width: 35,
                    height: 35,
                    fit: BoxFit.contain,
                  ),
            label: '首页',
          ),
          BottomNavigationBarItem(
            icon: _currentIndex == 1
                ? SvgPicture.asset(
                    'images/msg_selected.svg',
                    width: 35,
                    height: 35,
                    fit: BoxFit.contain,
                  )
                : SvgPicture.asset(
                    'images/msg_unselected.svg',
                    width: 35,
                    height: 35,
                    fit: BoxFit.contain,
                  ),
            label: '消息',
          ),
          // 如果用户没有登录的话，就展示默认的icon，否则展示用户头像
          avatar == ''
              ? BottomNavigationBarItem(
                  icon: _currentIndex == 2
                      ? SvgPicture.asset(
                          'images/me_selected.svg',
                          width: 40,
                          height: 40,
                          fit: BoxFit.contain,
                        )
                      : SvgPicture.asset(
                          'images/me_unselected.svg',
                          width: 40,
                          height: 40,
                          fit: BoxFit.contain,
                        ),
                  label: '我',
                )
              : BottomNavigationBarItem(
                  icon: ClipOval(
                      child: CircleAvatar(
                    radius: 24,
                    backgroundColor:
                        _currentIndex == 2 ? Colors.blue[700] : Colors.white,
                    child: CircleAvatar(
                      radius: 22, // 确保边框宽度
                      backgroundImage: NetworkImage(avatar),
                      backgroundColor: Colors.transparent,
                    ),
                  )),
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
