import 'package:flutter/material.dart';

class SaySomething extends StatelessWidget {
  final String avatar;
  // 传递用户的头像
  const SaySomething({super.key, required this.avatar});
  

  @override
  Widget build(BuildContext context) {
    return Row(
      children: <Widget>[
        // 左侧是我的头像
        CircleAvatar(
          backgroundImage: NetworkImage(avatar),
        ),
        // 右侧是输入框和发送按钮
        const Expanded(
          child: TextField(
            decoration: InputDecoration(
              hintText: '说点什么...',
            ),
          ),
        ),
        IconButton(
          icon: const Icon(Icons.send),
          onPressed: () {
            // TODO: 在这里添加发送消息的逻辑
          },
        ),
      ],
    );
  }
}