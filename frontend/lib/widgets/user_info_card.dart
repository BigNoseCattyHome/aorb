import 'package:flutter/material.dart';
import '../models/user.dart';  // 确保你有User模型

class UserInfoCard extends StatelessWidget {
  final User user;

  const UserInfoCard({super.key, required this.user});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        leading: CircleAvatar(backgroundImage: NetworkImage(user.avatarUrl)),
        title: Text(user.nickname),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: <Widget>[
            Text('关注: ${user.followed.length}'),
            Text('粉丝: ${user.followers.length}'),
            Text('金币: ${user.coins}'),
          ],
        ),
      ),
    );
  }
}
