import 'package:flutter/material.dart';

class PublicVote extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('公开投票'),
      ),
      body: Center(
        child: Text('公开展示你的投票，让大家帮你选择。'),
      ),
    );
  }
}
