import 'package:flutter/material.dart';

class MessageReply extends StatelessWidget {
  final String message;
  final String time;
  final String pollId;
  final String content;

  const MessageReply({
    Key? key,
    required this.message,
    required this.time,
    required this.pollId,
    required this.content,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      padding: const EdgeInsets.all(16.0),
      decoration: BoxDecoration(
        gradient: const LinearGradient(
          colors: [Colors.blue, Colors.purple],
        ),
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
          Text(message, style: const TextStyle(color: Colors.white)),
          const SizedBox(height: 8.0),
          Text(time, style: const TextStyle(color: Colors.grey)),
        ],
      ),
    );
  }
}
