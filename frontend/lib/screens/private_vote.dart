import 'package:flutter/material.dart';

class PrivateVote extends StatelessWidget {
  const PrivateVote({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('私密投票'),
      ),
      body: const Center(
        child: Text('私密地进行投票，只有被邀请者可见。'),
      ),
    );
  }
}
