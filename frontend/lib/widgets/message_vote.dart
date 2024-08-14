import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/generated/poll.pb.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/screens/poll_detail_page.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/utils/time.dart';
import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

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
  int commentsCount = 0;

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
      UserRequest()
        ..username = widget.username
        ..fields.addAll(['avatar', 'bgpic_pollcard', 'nickname']),
    )
        .then((response) {
      setState(() {
        avatar = response.user.avatar;
        bgpicPollcard = response.user.bgpicPollcard;
        nickname = response.user.nickname;
      });
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
      });
    });
  }

  @override
  Widget build(BuildContext context) {
    return Slidable(
      endActionPane: ActionPane(
        motion: const ScrollMotion(),
        children: [
          SlidableAction(
            onPressed: (context) {
              widget.onRead();
            },
            backgroundColor: Colors.blue,
            foregroundColor: Colors.white,
            icon: Icons.check,
            label: '已读',
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
        child: Card(
          margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
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
                        style: const TextStyle(fontWeight: FontWeight.bold)),
                    const SizedBox(width: 8),
                    const Text("选择了"),
                    const SizedBox(width: 4),
                    Text(widget.choice,
                        style: const TextStyle(fontWeight: FontWeight.bold)),
                  ],
                ),
                const SizedBox(height: 16),
                Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Icon(Icons.chat_bubble_outline, size: 24),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Text(
                            title,
                            style: const TextStyle(
                                fontWeight: FontWeight.bold, fontSize: 18),
                          ),
                          const SizedBox(height: 4),
                          Text(content, style: const TextStyle(fontSize: 14)),
                        ],
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                Align(
                  alignment: Alignment.bottomRight,
                  child: Text(
                    formatTimestamp(widget.time, ""),
                    style: const TextStyle(fontSize: 12, color: Colors.grey),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
