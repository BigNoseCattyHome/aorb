import 'package:flutter/material.dart';
import '../widgets/user_info_card.dart';
import '../models/user.dart';

class Profile extends StatelessWidget {
  final User user = User(
    id: '1',
    nickname: '张三',
    avatarUrl: 'https://example.com/avatar.jpg',
    followed: ['2', '3'],
    followers: ['2', '3'],
    blacklist: [],
    coins: 100,
    coinsRecord: ['+10', '-5'],
    questionsAsk: ['1', '2'],
    questionsAnswer: ['1', '2'],
    channels: ['1', '2'],
  );

  Profile({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('我的'),
      ),
      body: ListView(
        children: <Widget>[
          UserInfoCard(user: user),
          ListTile(
            leading: const Icon(Icons.history),
            title: const Text('发起的投票历史记录'),
            onTap: () {},
          ),
          ListTile(
            leading: const Icon(Icons.history),
            title: const Text('参与的投票历史记录'),
            onTap: () {},
          ),
          ListTile(
            leading: const Icon(Icons.info_outline),
            title: const Text('关于开发者'),
            onTap: () {},
          ),
        ],
      ),
    );
  }
}
