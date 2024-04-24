import 'package:flutter/material.dart';
import '../models/vote.dart';  // 假设你已经定义了Vote模型
import '../widgets/choice_card.dart';

class HelpThemChoose extends StatelessWidget {
  // 这里只是示意，实际上你可能需要从服务器获取数据
  final List<Vote> votes = [];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('帮他选'),
      ),
      body: ListView.builder(
        itemCount: votes.length,
        itemBuilder: (context, index) => ChoiceCard(vote: votes[index]),
      ),
    );
  }
}
