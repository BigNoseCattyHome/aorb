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

class MessageVote extends StatefulWidget {
  final String username; // 需要根据 username 获取头像,昵称
  final String pollId; // 需要根据 pollId 获取title和content
  final Timestamp time; // 消息时间
  final String choice; // 他人的投票选择
  final Function onRead; // 标记为已读的回调函数

  const MessageVote({
    Key? key,
    required this.username,
    required this.time,
    required this.choice,
    required this.pollId,
    required this.onRead,
  }) : super(key: key);

  @override
  MessageVoteState createState() => MessageVoteState();
}

class MessageVoteState extends State<MessageVote> {
  String avatar = '';
  String bgpicPollcard = '';
  String nickname = '';
  String title = '';
  String content = '';
  String pollowner = '';
  int commentsCount = 0;
  Color _textColor = Colors.black; // 默认文字颜色

  @override
  void initState() {
    super.initState();
    fetchUserInfo();
    fetchPollInfo();
  }

  // 根据username查询avatar和bgpic_pollcard
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
        pollowner = response.poll.username;
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
                username: pollowner,
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
                              fontWeight: FontWeight.bold, color: _textColor)),
                      const SizedBox(width: 8),
                      Text(
                        "选择了",
                        style: TextStyle(color: _textColor),
                      ),
                      const SizedBox(width: 4),
                      Text(widget.choice,
                          style: TextStyle(
                              fontWeight: FontWeight.bold, color: _textColor)),
                    ],
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
                                style:
                                    TextStyle(fontSize: 14, color: _textColor)),
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
        ),
      ),
    );
  }
}
