import 'package:flutter/material.dart';
import '../models/vote.dart';  // 确保你有Vote模型

class ChoiceCard extends StatelessWidget {
  final Vote vote;

  const ChoiceCard({super.key, required this.vote});

  @override
  Widget build(BuildContext context) {
    return Card(
      child: ListTile(
        title: Text(vote.title),
        subtitle: Text(vote.description),
        trailing: const Icon(Icons.arrow_forward),
        onTap: () {
          // 处理点击事件
        },
      ),
    );
  }
}
