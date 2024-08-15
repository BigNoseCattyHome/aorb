import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/screens/user_profile_page.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/utils/container_decoration.dart';
import 'package:aorb/utils/time.dart';
import 'package:aorb/utils/color_analyzer.dart';
import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

class MessageFollowed extends StatefulWidget {
  final String username;
  final Timestamp time;
  final Function onRead; // 标记为已读的回调函数

  const MessageFollowed({
    Key? key,
    required this.username,
    required this.time,
    required this.onRead,
  }) : super(key: key);

  @override
  MessageFollowedState createState() => MessageFollowedState();
}

class MessageFollowedState extends State<MessageFollowed> {
  String avatar = '';
  String bgpicPollcard = '';
  String nickname = '';
  Color _textColor = Colors.black; // 默认文字颜色

  @override
  void initState() {
    super.initState();
    fetchUserInfo();
  }

  // 获取用户信息并设置文字颜色
  void fetchUserInfo() {
    UserService()
        .getUserInfo(
      UserRequest()..username = widget.username,
    )
        .then((response) async {
      setState(() {
        avatar = response.user.avatar;
        bgpicPollcard = response.user.bgpicPollcard;
        nickname = response.user.nickname;
      });
      // 使用 ColorAnalyzer 获取适合的文字颜色
      _textColor = await ColorAnalyzer.getTextColor(bgpicPollcard);
      if (mounted) {
        setState(() {});
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Slidable(
      endActionPane: ActionPane(
        motion: const ScrollMotion(),
        extentRatio: 0.2, // 控制滑动的范围
        children: [
          SlidableAction(
            onPressed: (context) {
              widget.onRead();
            },
            backgroundColor: Colors.blue,
            foregroundColor: Colors.white,
            icon: Icons.check,
            label: '已读',
            borderRadius: BorderRadius.circular(10),
          ),
        ],
      ),
      child: GestureDetector(
          onTap: () {
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) =>
                    UserProfilePage(username: widget.username),
              ),
            );
          },
          child: Container(
            decoration: createBackgroundDecoration(bgpicPollcard),
            margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            child: Card(
              color: Colors.transparent,
              shadowColor: Colors.transparent,
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        CircleAvatar(
                          backgroundImage: NetworkImage(avatar),
                          radius: 20,
                        ),
                        const SizedBox(width: 8),
                        Text(nickname,
                            style: TextStyle(
                                fontWeight: FontWeight.bold,
                                color: _textColor)),
                        const SizedBox(width: 8),
                        Text("关注了你",
                            style: TextStyle(fontSize: 16, color: _textColor)),
                      ],
                    ),
                    Align(
                      alignment: Alignment.bottomRight,
                      child: Text(
                        formatTimestamp(widget.time, ""),
                        style: TextStyle(
                            fontSize: 12, color: _textColor.withOpacity(0.6)),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          )),
    );
  }
}
