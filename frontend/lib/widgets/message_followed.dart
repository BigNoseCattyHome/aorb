import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/utils/time.dart';
import 'package:flutter/material.dart';

class MessageFollowed extends StatelessWidget {
  final String username;
  final Timestamp time;

  const MessageFollowed({
    Key? key,
    required this.username,
    required this.time,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      padding: const EdgeInsets.all(16.0),
      decoration: BoxDecoration(
        color: Colors.blue.shade800,
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
          CircleAvatar(
            radius: 20.0,
            backgroundColor: Colors.white,
            child: Text(
              username[0].toUpperCase(),
              style: TextStyle(color: Colors.blue.shade800),
            ),
          ),
          Text(username, style: const TextStyle(color: Colors.white)),
          const SizedBox(height: 8.0),
          Text(formatTimestamp(time),
              style: const TextStyle(color: Colors.grey)),
        ],
      ),
    );
  }
}
