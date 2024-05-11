import 'package:flutter/material.dart';
import '../models/comment.dart';  // 确保你有Comment模型

class CommentSection extends StatelessWidget {
  final List<Comment> comments;

  const CommentSection({super.key, required this.comments});

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: comments.length,
      itemBuilder: (context, index) {
        final comment = comments[index];
        return ListTile(
          // ! 这里记得把头像的URL替换成真实的URL
          leading: CircleAvatar(backgroundImage: NetworkImage(comment.userid)),
          title: Text(comment.userid),
          subtitle: Text(comment.advise),
        );
      },
    );
  }
}
