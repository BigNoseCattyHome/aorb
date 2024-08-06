// home_page.dart
import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:flutter/material.dart';
import 'package:aorb/widgets/poll_card.dart'; // 引入未投票问题组件
import 'dart:async';

class HomePage extends StatefulWidget {
  final TabController tabController;

  const HomePage({super.key, required this.tabController});

  @override
  HomePageState createState() => HomePageState();
}

class HomePageState extends State<HomePage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController; // 顶部导航栏控制器
  late Future<List<PollCard>> _futureQuestions;

  @override
  void initState() {
    super.initState(); // 调用父类的 initState 方法
    _tabController = widget.tabController; // 初始化顶部导航栏控制器
    _futureQuestions = _fetchQuestions(); // 初始化时调用 _fetchQuestions 获取数据
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // 顶部栏，由main_page接管

      // 发布按钮
      floatingActionButton: FloatingActionButton(
          onPressed: () {
            // 导航到发布界面
          },
          backgroundColor: Colors.blue[700],
          child: const Icon(
            Icons.add,
            color: Colors.white,
            size: 30,
          )),

      // 中间的投票卡片
      body: TabBarView(controller: _tabController, children: [
        FutureBuilder<List<PollCard>>(
          future: _futureQuestions,
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return const Center(child: CircularProgressIndicator());
            } else if (snapshot.hasError) {
              return Center(child: Text('Error: ${snapshot.error}'));
            } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
              return const Center(child: Text('No questions available.'));
            } else {
              return ListView.builder(
                itemCount: snapshot.data!.length,
                itemBuilder: (context, index) {
                  return snapshot.data![index];
                },
              );
            }
          },
        ),
        FutureBuilder<List<PollCard>>(
          future: _futureQuestions,
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return const Center(child: CircularProgressIndicator());
            } else if (snapshot.hasError) {
              return Center(child: Text('Error: ${snapshot.error}'));
            } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
              return const Center(child: Text('No questions available.'));
            } else {
              return ListView.builder(
                itemCount: snapshot.data!.length,
                itemBuilder: (context, index) {
                  return Padding(
                    padding: const EdgeInsets.symmetric(
                        horizontal: 16.0, vertical: 8.0), // 设置外边距
                    child: snapshot.data![index],
                  );
                },
              );
            }
          },
        ),
      ]),
    );
  }

  Future<List<PollCard>> _fetchQuestions() async {
    // 模拟从服务器获取数据
    await Future.delayed(const Duration(seconds: 2));
    return [
      PollCard(
        title: "午饭吃什么呀?",
        content: '想了半天没有想出来到底要吃什么，好纠结，真可恶！',
        options: const ['麻辣烫', '炸鸡汉堡'],
        votePercentage: const [0.4, 0.6],
        voteCount: 20,
        time: Timestamp.fromDateTime(DateTime.parse("2024-07-01 12:05:00")),
        avatar: 'https://s2.loli.net/2024/05/27/2MgJcvLtOVKmAdn.jpg',
        nickname: '爱吃饭的小袁同学',
        userId: '1',
        pollId: '1',
        backgroundImage: 'gradient:0xFF8DB3EB,0xFFF895CA',
        selectedOption: -1,
      ),
      PollCard(
        title: '下午去哪里玩?',
        content: '天气很好，感觉顾村公园和外滩都挺不错的，选哪个？',
        options: const ['顾村公园', '外滩'],
        votePercentage: const [0.16, 0.84],
        voteCount: 76,
        time: Timestamp.fromDateTime(DateTime.parse("2024-07-01 12:05:00")),
        avatar: 'https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png',
        nickname: '花枝鼠gogo来帮忙',
        userId: '2',
        pollId: '2',
        backgroundImage: '0xFF354967',
        selectedOption: -1,
      ),
      PollCard(
        title: '要不要去小美家玩啊？',
        content: '小美小美小美好香的小美，Bad 小新',
        options: const ['麻辣烫', '炸鸡汉堡'],
        votePercentage: const [0.4, 0.6],
        voteCount: 20,
        time: Timestamp.fromDateTime(DateTime.parse("2024-07-01 12:05:00")),
        avatar: 'https://s2.loli.net/2024/05/27/QzKM41C3Vs5FeHW.jpg',
        nickname: '风见澈Siri',
        userId: '3',
        pollId: '3',
        backgroundImage: 'https://s2.loli.net/2024/05/25/HqJM8dTuSRbUNBO.jpg',
        selectedOption: -1,
      ),
      PollCard(
        title: 'Exploring the Enigmatic World of Quantum Mechanics',
        content:
            'Quantum mechanics is a fundamental theory in physics that provides a description of the physical properties of nature at the scale of atoms and subatomic particles. However, it\'s not as straightforward as classical physics. Can you help me understand some of the key concepts and phenomena of quantum mechanics?',
        options: const [
          'Yes, I\'d love to.',
          'I\'m familiar with quantum mechanics.',
        ],
        votePercentage: const [0.4, 0.6],
        voteCount: 20,
        time: Timestamp.fromDateTime(DateTime.parse("2024-07-01 12:05:00")),
        avatar: 'https://s2.loli.net/2024/05/27/alt3BKPYhzmV4E7.jpg',
        nickname: 'Anti Cris',
        userId: '4',
        pollId: '4',
        backgroundImage: 'gradient:0x7FFCE300,0xFFFF5065,0xFF1F9AC1',
        selectedOption: -1,
      ),
    ];
  }
}
