// home_page.dart
import 'package:flutter/material.dart';
import '../widgets/top_bar_index.dart'; // 引入顶部导航栏组件
import '../widgets/question_unvoted.dart'; // 引入未投票问题组件
import 'dart:async';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  HomePageState createState() => HomePageState();
}

class HomePageState extends State<HomePage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController; // 顶部导航栏控制器
  late Future<List<QuestionUnvoted>> _futureQuestions;

  @override
  void initState() {
    super.initState(); // 调用父类的 initState 方法
    _tabController = TabController(length: 3, vsync: this); // 初始化顶部导航栏控制器
    _futureQuestions = _fetchQuestions(); // 初始化时调用 _fetchQuestions 获取数据
  }

  Future<List<QuestionUnvoted>> _fetchQuestions() async {
    // 模拟从服务器获取数据
    await Future.delayed(const Duration(seconds: 2));
    return [
      const QuestionUnvoted(
        title: '午饭吃什么呀?',
        content: '想了半天没有想出来到底要吃什么，好纠结，真可恶！',
        options: ['麻辣烫', '炸鸡汉堡'],
        votePercentage: [0.4, 0.6],
        voteCount: 20,
        time: '发布于 5 小时前',
        avatar: 'https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png',
        nickname: '爱吃饭的小袁同学',
        questionId: '1',
        backgroundImage: 'https://s2.loli.net/2024/05/25/HqJM8dTuSRbUNBO.jpg',
        selectedOption: -1,
      ),
      const QuestionUnvoted(
        title: '午饭吃什么呀?',
        content: '想了半天没有想出来到底要吃什么，好纠结，真可恶！',
        options: ['麻辣烫', '炸鸡汉堡'],
        votePercentage: [0.4, 0.6],
        voteCount: 20,
        time: '发布于 5 小时前',
        avatar: 'https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png',
        nickname: '爱吃饭的小袁同学',
        questionId: '1',
        backgroundImage: 'https://s2.loli.net/2024/05/25/HqJM8dTuSRbUNBO.jpg',
        selectedOption: -1,
      ),
      const QuestionUnvoted(
        title: '午饭吃什么呀?',
        content: '想了半天没有想出来到底要吃什么，好纠结，真可恶！',
        options: ['麻辣烫', '炸鸡汉堡'],
        votePercentage: [0.4, 0.6],
        voteCount: 20,
        time: '发布于 5 小时前',
        avatar: 'https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png',
        nickname: '爱吃饭的小袁同学',
        questionId: '1',
        backgroundImage: 'https://s2.loli.net/2024/05/25/HqJM8dTuSRbUNBO.jpg',
        selectedOption: -1,
      ),
      const QuestionUnvoted(
        title: '午饭吃什么呀?',
        content: '想了半天没有想出来到底要吃什么，好纠结，真可恶！',
        options: ['麻辣烫', '炸鸡汉堡'],
        votePercentage: [0.4, 0.6],
        voteCount: 20,
        time: '发布于 5 小时前',
        avatar: 'https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png',
        nickname: '爱吃饭的小袁同学',
        questionId: '1',
        backgroundImage: 'https://s2.loli.net/2024/05/25/HqJM8dTuSRbUNBO.jpg',
        selectedOption: -1,
      ),
      const QuestionUnvoted(
        title: '午饭吃什么呀?',
        content: '想了半天没有想出来到底要吃什么，好纠结，真可恶！',
        options: ['麻辣烫', '炸鸡汉堡'],
        votePercentage: [0.4, 0.6],
        voteCount: 20,
        time: '发布于 5 小时前',
        avatar: 'https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png',
        nickname: '爱吃饭的小袁同学',
        questionId: '1',
        backgroundImage: 'https://s2.loli.net/2024/05/25/HqJM8dTuSRbUNBO.jpg',
        selectedOption: -1,
      ),
      // 可以添加更多QuestionUnvoted示例
    ];
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
      appBar: const DynamicTopBar(
        tabs: ['推荐', '关注'],
        showSearch: true,
      ),

      // 中间的投票卡片
      body: 
      FutureBuilder<List<QuestionUnvoted>>(
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

    );
  }
}
