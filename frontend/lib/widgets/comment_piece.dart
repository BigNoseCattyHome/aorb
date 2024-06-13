import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class CommentPiece extends StatelessWidget {
  final String avatar;
  final String nickname;
  final String content;
  final String ipdress;
  final DateTime time;

  const CommentPiece(
      {super.key,
      required this.avatar,
      required this.content,
      required this.ipdress,
      required this.nickname,
      required this.time});

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            CircleAvatar(
              backgroundImage: NetworkImage(avatar),
            ),
            SizedBox(width: 8),
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  nickname,
                  style: TextStyle(fontWeight: FontWeight.bold),
                ),
                Text(
                  ipdress,
                  style: TextStyle(color: Colors.grey, fontSize: 12),
                ),
              ],
            ),
          ],
        ),
        SizedBox(height: 8),
        Text(content),
        SizedBox(height: 8),
        Text(
          DateFormat('yyyy-MM-dd HH:mm').format(time),
          style: TextStyle(color: Colors.grey, fontSize: 12),
        ),
      ],
    );
  }
}
