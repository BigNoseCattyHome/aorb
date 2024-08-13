import 'package:flutter/material.dart';
import 'package:aorb/services/comment_service.dart';
import 'package:aorb/generated/comment.pbgrpc.dart';

class SaySomething extends StatelessWidget {
  final String currentUserAvatar; // 传递用户的头像
  final String currentUsername; // 当前用户的username
  final String pollId; // 当前投票的id
  final VoidCallback onCommentPosted;

  SaySomething(
      {super.key,
      required this.currentUsername,
      required this.currentUserAvatar,
      required this.pollId,
      required this.onCommentPosted});

  final _textController = TextEditingController(); // 获取输入框中的内容

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
      child: Row(
        children: <Widget>[
          // 左侧是我的头像
          CircleAvatar(
            backgroundImage: NetworkImage(currentUserAvatar),
            radius: 20.0, // 调整头像大小
          ),
          const SizedBox(width: 16.0), // 增加头像和输入框之间的间距
          // 右侧是输入框和发送按钮
          Expanded(
            child: TextField(
              controller: _textController,
              decoration: InputDecoration(
                hintText: '说点什么...',
                contentPadding: const EdgeInsets.symmetric(
                    horizontal: 16.0, vertical: 12.0), // 调整输入框内边距
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24.0), // 圆角边框
                ),
              ),
            ),
          ),
          const SizedBox(width: 5.0), // 增加输入框和发送按钮之间的间距
          IconButton(
            icon: const Icon(Icons.send),
            color: Colors.blue, // 设置发送按钮为蓝色
            onPressed: () {
              String inputText = _textController.text;
              _textController.clear(); // 点击发送按钮后，清空输入框
              CommentService().actionComment(currentUsername, pollId,
                  ActionCommentType.ACTION_COMMENT_TYPE_ADD, inputText);
              onCommentPosted();
            },
          ),
        ],
      ),
    );
  }
}
