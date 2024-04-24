import 'package:flutter/material.dart';

class SystemChoice extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('系统选择'),
      ),
      body: Center(
        child: Text('系统自动帮你做出选择。'),
      ),
    );
  }
}
