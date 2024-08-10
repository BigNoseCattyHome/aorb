import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/utils/time.dart';
import 'package:flutter/material.dart';

class MessageVote extends StatelessWidget {
  final String username; // 需要根据 username 获取头像
  final String pollId; // 需要根据 pollId 获取title和content
  final Timestamp time; // 消息时间
  final String choice; // 他人的投票选择
  final String userId; // 暂时不用

  const MessageVote({
    Key? key,
    required this.username,
    required this.userId,
    required this.time,
    required this.choice,
    required this.pollId,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      padding: const EdgeInsets.all(16.0),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.0),
        boxShadow: const [
          BoxShadow(
            color: Colors.black26,
            blurRadius: 6.0,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // 投票人，选择

          // 问题

          // 时间
          Text(formatTimestamp(time),
              style: const TextStyle(color: Colors.grey)),
        ],
      ),
    );
  }
}
