// messages_page.dart
import 'package:flutter/material.dart';
import '../widgets/top_bar_index.dart'; // 引入顶部导航栏组件

class MessagesPage extends StatefulWidget {
  const MessagesPage({super.key});

  @override
  MessagesPageState createState() => MessagesPageState();
}

class MessagesPageState extends State<MessagesPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController; // 顶部导航栏控制器

  @override
  void initState() {
    super.initState(); // 调用父类的 initState 方法
    _tabController = TabController(length: 2, vsync: this); // 初始化顶部导航栏控制器
  }

  @override
  void dispose() {
    // 释放资源
    _tabController.dispose();
    super.dispose();
  }
  @override
  Widget build(BuildContext context) {
    return Scaffold(
        // 顶部栏
        appBar: DynamicTopBar(
      tabs: const ['提醒', '私信'],
      showSearch: true,
      tabController: _tabController,
    ));
  }
}
