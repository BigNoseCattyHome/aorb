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
  Widget build(BuildContext context) {
    return const Scaffold(
      // 顶部栏
      appBar: DynamicTopBar(
        tabs: ['提醒', '私信'],
        showSearch: true,
      )


    );
  }
}
