import 'package:flutter/material.dart';

class AnonymousVote extends StatelessWidget {
  const AnonymousVote({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('匿名投票'),
      ),
      body: const Center(
        child: Text('在这里参与匿名投票，保护您的隐私'),
      ),
    );
  }
}
