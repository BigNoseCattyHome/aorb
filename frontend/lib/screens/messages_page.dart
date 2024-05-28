// messages_page.dart
import 'package:flutter/material.dart';

class MessagesPage extends StatefulWidget {
  final TabController tabController;

  const MessagesPage({super.key, required this.tabController});

  @override
  MessagesPageState createState() => MessagesPageState();
}

class MessagesPageState extends State<MessagesPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController; // 顶部导航栏控制器

  @override
  void initState() {
    super.initState(); // 调用父类的 initState 方法
    _tabController = widget.tabController; // 初始化顶部导航栏控制器
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return const Scaffold(

        );
  }
}
