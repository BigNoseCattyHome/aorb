import 'package:flutter/material.dart';

class AnonymousVote extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('匿名投票'),
      ),
      body: Center(
        child: Text('在这里参与匿名投票，保护您的隐私'),
      ),
    );
  }
}
