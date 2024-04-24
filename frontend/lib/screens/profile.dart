import 'package:flutter/material.dart';
import '../models/user.dart';  // 假设你已经定义了User模型
import '../widgets/user_info_card.dart';
import '../models/user.dart';

class Profile extends StatelessWidget {
  final User user = User(
    id: '001',
    nickname: 'JohnDoe',
    avatarUrl: 'https://via.placeholder.com/150',
    following: ['002', '003'],
    followers: ['004', '005', '006'],
    coins: 150,
    coinTransactions:['007', '008', '009', '010', '011', '012', '013', '014'],
    blocked: [],
  );

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('我的'),
      ),
      body: ListView(
        children: <Widget>[
          UserInfoCard(user: user),
          ListTile(
            leading: Icon(Icons.history),
            title: Text('发起的投票历史记录'),
            onTap: () {},
          ),
          ListTile(
            leading: Icon(Icons.history),
            title: Text('参与的投票历史记录'),
            onTap: () {},
          ),
          ListTile(
            leading: Icon(Icons.info_outline),
            title: Text('关于开发者'),
            onTap: () {},
          ),
        ],
      ),
    );
  }
}
