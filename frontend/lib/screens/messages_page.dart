import 'package:flutter/material.dart';

// 定义 MessagesPage 组件，它是一个有状态的组件（StatefulWidget）
class MessagesPage extends StatefulWidget {
  // 传入的 TabController，用于控制顶部导航栏的切换
  final TabController tabController;

  // 构造函数
  const MessagesPage({super.key, required this.tabController});

  // 创建状态类
  @override
  MessagesPageState createState() => MessagesPageState();
}

// 状态类，管理 MessagesPage 的状态
class MessagesPageState extends State<MessagesPage> {
  // 顶部导航栏控制器
  late TabController _tabController;

  // 初始化状态
  @override
  void initState() {
    super.initState();
    _tabController =
        widget.tabController; // 将传入的 TabController 赋值给 _tabController
  }

  // 释放资源
  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  // 构建页面 UI
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // appBar 由 MainPage 管理

      // body 部分通过 TabBarView 显示不同的页面
      body: TabBarView(
        controller: _tabController, // 使用 _tabController 控制页面切换
        children: [
          // 第一个页面：提醒
          _buildNoticeList(),
          // 第二个页面：私信
          _buildMessageList()
        ],
      ),
    );
  }

  // 构建消息列表
  Widget _buildNoticeList() {
    // 示例消息数据
    final messages = [];

    // 使用 ListView.builder 动态生成列表项
    return ListView.builder(
      itemCount: messages.length, // 列表项数量
      itemBuilder: (context, index) {
        final message = messages[index];
        return ListTile(
          title: Text(message['title'] as String), // 消息标题
          subtitle: Text(message['content'] as String), // 消息内容
          // 如果消息未读，则显示红点图标，否则不显示
          trailing: message['isRead'] as bool
              ? null
              : const Icon(Icons.circle, color: Colors.red, size: 10),
          // 点击消息项，显示消息详情
          onTap: () {
            _showMessageDetail(message as Map<String, String>);
          },
        );
      },
    );
  }

  Widget _buildMessageList() {
    return Container();
  }

  // 显示消息详情页面
  void _showMessageDetail(Map<String, String> message) {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => MessageDetailPage(message: message), // 创建消息详情页面
      ),
    );
  }
}

// 消息详情页面
class MessageDetailPage extends StatelessWidget {
  // 传入的消息数据
  final Map<String, String> message;

  // 构造函数
  const MessageDetailPage({super.key, required this.message});

  // 构建消息详情页面 UI
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(message['title']!), // 标题栏，显示消息标题
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Text(message['content']!), // 显示消息内容
      ),
    );
  }
}
