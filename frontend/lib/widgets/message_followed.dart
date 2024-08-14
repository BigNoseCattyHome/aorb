import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/screens/user_profile_page.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/utils/container_decoration.dart';
import 'package:aorb/utils/time.dart';
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

  @override
  void initState() {
    super.initState();
    fetchUserInfo();
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
                            style:
                                const TextStyle(fontWeight: FontWeight.bold)),
                        const SizedBox(width: 8),
                        const Text("关注了你", style: TextStyle(fontSize: 16)),
                      ],
                    ),
                    const SizedBox(height: 16),
                    Align(
                      alignment: Alignment.bottomRight,
                      child: Text(
                        formatTimestamp(widget.time, ""),
                        style:
                            const TextStyle(fontSize: 12, color: Colors.grey),
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
