import 'package:flutter/material.dart';

class PublicVote extends StatelessWidget {
  const PublicVote({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('公开投票'),
      ),
      body: const Center(
        child: Text('公开展示你的投票，让大家帮你选择。'),
      ),
    );
  }
}
