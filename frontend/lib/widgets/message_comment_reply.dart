import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/generated/poll.pb.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/screens/poll_detail_page.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/utils/color_analyzer.dart';
import 'package:aorb/utils/container_decoration.dart';
import 'package:aorb/utils/time.dart';
import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';
import 'package:flutter_svg/flutter_svg.dart';

class MessageCommentReply extends StatefulWidget {
  final String username;
  final String replyText;
  final String pollId;
  final Timestamp time;
  final Function onRead; // 标记为已读的回调函数

  const MessageCommentReply({
    Key? key,
    required this.username,
    required this.replyText,
    required this.pollId,
    required this.time,
    required this.onRead,
  }) : super(key: key);

  @override
  MessageCommentReplyState createState() => MessageCommentReplyState();
}

class MessageCommentReplyState extends State<MessageCommentReply> {
  String avatar = '';
  String bgpicPollcard = '';
  String nickname = '';
  String title = '';
  String content = '';
  int commentsCount = 0;
  Color _textColor = Colors.black; // 默认文字颜色

  @override
  void initState() {
    super.initState();
    fetchUserInfo();
    fetchPollInfo();
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

  // 获取投票信息
  void fetchPollInfo() {
    PollService()
        .getPoll(
      GetPollRequest()..pollUuid = widget.pollId,
    )
        .then((response) {
      setState(() {
        title = response.poll.title;
        content = response.poll.content;
        commentsCount = response.poll.commentList.length;
      });
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
                builder: (context) => PollDetailPage(
                  pollId: widget.pollId,
                  username: widget.username,
                  userId: '',
                ),
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
                        Text("回复了你",
                            style: TextStyle(
                                fontSize: 16,
                                fontWeight: FontWeight.normal,
                                color: _textColor)),
                      ],
                    ),
                    const SizedBox(height: 16),
                    Center(
                      child: Text(
                        widget.replyText,
                        style: TextStyle(
                          fontSize: 28,
                          fontWeight: FontWeight.bold,
                          color: _textColor,
                        ),
                        textAlign: TextAlign.center,
                      ),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        SvgPicture.asset('images/comments.svg'),
                        const SizedBox(width: 8),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                title,
                                style: TextStyle(
                                    fontWeight: FontWeight.bold,
                                    fontSize: 18,
                                    color: _textColor),
                              ),
                              const SizedBox(height: 4),
                              Text(content,
                                  style: TextStyle(
                                      fontSize: 14, color: _textColor)),
                            ],
                          ),
                        ),
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
