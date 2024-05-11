import 'package:flutter/material.dart';

class SystemChoice extends StatelessWidget {
  const SystemChoice({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('系统选择'),
      ),
      body: const Center(
        child: Text('系统自动帮你做出选择。'),
      ),
    );
  }
}
