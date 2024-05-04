import 'package:flutter/material.dart';
import '../models/vote.dart';
import '../widgets/choice_card.dart';
// import '../widgets/customed_app_bar.dart';

class Index extends StatefulWidget {
  @override
  _IndexState createState() => _IndexState();
}

class _IndexState extends State<Index> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final List<Vote> votes = [
    Vote(
      id: '1',
      type: 'single',
      time: '2021-07-01 12:00:00',
      sponsor: '张三',
      title: '你喜欢吃什么水果？',
      description: '这是一个关于水果的投票',
      options: ['苹果', '香蕉', '橙子'],
      channel: '1',
      comments: [],
      fee: 0,
      inviteIDs: [],
      voters: ['2', '3'],
    ),
    Vote(
      id: '2',
      type: 'multiple',
      time: '2021-07-02 12:00:00',
      sponsor: '李四',
      title: '你喜欢哪些颜色？',
      description: '这是一个关于颜色的投票',
      options: ['红色', '黄色', '蓝色'],
      channel: '2',
      comments: [],
      fee: 0,
      inviteIDs: [],
      voters: ['1', '3'],
    ),
  ]; // 模拟数据

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: TabBar(
          controller: _tabController,
          indicator: UnderlineTabIndicator(
            borderSide: BorderSide(width: 2.0, color: Colors.blue),
          ),
          labelColor: Colors.black,
          tabs: [
            Tab(text: '动态'),
            Tab(text: '推荐'),
          ],
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          ListView.builder(
            itemCount: votes.length,
            itemBuilder: (context, index) => ChoiceCard(vote: votes[index]),
          ),
          ListView.builder(
            itemCount: votes.length,
            itemBuilder: (context, index) => ChoiceCard(vote: votes[index]),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.pushNamed(context, '/help_me_choose');
        },
        child: Icon(Icons.add),
        backgroundColor: Colors.blue,
      ),
    );
  }
}
