import 'package:flutter/material.dart';
import '../models/comment.dart';  // 确保你有Comment模型

class CommentSection extends StatelessWidget {
  final List<Comment> comments;

  CommentSection({required this.comments});

  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: comments.length,
      itemBuilder: (context, index) {
        final comment = comments[index];
        return ListTile(
          leading: CircleAvatar(backgroundImage: NetworkImage(comment.userAvatarUrl)),
          title: Text(comment.userName),
          subtitle: Text(comment.content),
        );
      },
    );
  }
}
