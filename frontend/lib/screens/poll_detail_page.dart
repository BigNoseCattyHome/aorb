import 'package:aorb/generated/poll.pb.dart';
import 'package:flutter/material.dart';

import 'package:aorb/services/user_service.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/widgets/say_something.dart';
import 'package:aorb/widgets/comment_piece.dart';
import 'package:aorb/widgets/poll_detail.dart';
import 'package:aorb/generated/user.pbgrpc.dart';

class PollDetailPage extends StatefulWidget {
  final String userId; // 用户Id，用于获取和博主的关注状态和用户信息
  final String pollId; // 这篇帖子的uuid，用于获取帖子详情
  const PollDetailPage({super.key, required this.pollId, required this.userId});

  @override
  PollDetailPageState createState() => PollDetailPageState();
}

class PollDetailPageState extends State<PollDetailPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  late bool isFollowed; // 是否关注
  late int cntComments; // 评论数
  late Poll poll; // 评论的具体内容
  late User user; // 用户信息

  @override
  void initState() {
    super.initState();

    // 获取关注状态
    IsUserFollowingRequest requestFollow = IsUserFollowingRequest()
      ..username = widget.userId;
    UserService().isUserFollowing(requestFollow).then((isFollowed) {
      setState(() {
        this.isFollowed = isFollowed;
      });
    });

    // 获取用户信息: nickname, avatar
    UserRequest request = UserRequest()..username = widget.userId;
    UserService().getUserInfo(request).then((response) {
      setState(() {
        user = response.user;
      });
    });

    // 获取投票详情
    PollService()
        .GetPoll(GetPollRequest()..pollUuid = widget.pollId)
        .then((pollResponse) {
      setState(() {
        poll = pollResponse.poll;
        cntComments = pollResponse.poll.commentList.length;
      });
    });

    _tabController = TabController(length: 3, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // 顶部栏，包括返回键，用户头像和用户名，关注按钮
      appBar: AppBar(
          // 返回键
          leading: IconButton(
            icon: const Icon(Icons.arrow_back),
            onPressed: () {
              Navigator.pop(context);
            },
          ),

          // 用户头像和用户名
          title: const Row(
            children: [
              CircleAvatar(
                // backgroundImage: NetworkImage(user.avatar),
                // ~ 测试用的图片
                backgroundImage: NetworkImage(
                    'https://s2.loli.net/2024/05/27/2MgJcvLtOVKmAdn.jpg'),
                radius: 15,
              ),
              SizedBox(width: 5),
              Text('爱吃饭的小袁同学'),
            ],
          ),

          // 椭圆矩形关注按钮，查询是否关注过
          // 没关注过就是蓝框透明的关注按钮；关注过了就是灰色按钮“已关注”
          actions: isFollowed
              ? [
                  // 蓝框透明的椭圆矩形关注按钮
                  TextButton(
                    style: TextButton.styleFrom(
                      side: const BorderSide(color: Colors.blue),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(20),
                      ),
                    ),
                    onPressed: () {
                      // TODO 在这里添加取消关注的逻辑
                    },
                    child:
                        const Text('已关注', style: TextStyle(color: Colors.blue)),
                  ),
                ]
              : [
                  TextButton(
                    style: TextButton.styleFrom(
                      side: const BorderSide(color: Colors.blue),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(20),
                      ),
                    ),
                    onPressed: () {
                      // TODO 在这里添加关注的逻辑
                    },
                    child:
                        const Text('关注', style: TextStyle(color: Colors.blue)),
                  ),
                ]),

      body: Column(children: [
        // 内容详情
        PollDetail(
            title: poll.title,
            content: poll.content,
            options: poll.options,
            votePercentage: poll.optionsCount.map((value) {
              return poll.optionsCount.reduce((a, b) => a + b) > 0
                  ? (value / poll.optionsCount.reduce((a, b) => a + b)) * 100
                  : 0.0;
            }).toList(),
            time: poll.createAt.toDateTime(),
            ipaddress: poll.ipaddress),

        // 分割线
        const Divider(),

        // 评论区
        Text('共有$cntComments条评论'),
        const SaySomething(avatar: 'avatar'),

        // TODO 从poll.comments中获取评论内容
        CommentPiece(
          avatar: 'https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png',
          nickname: '花枝鼠来帮忙',
          content:
              '这件事闹得挺大的，全球都传疯了。这件事闹得确实挺大，但也不是特别大，你要说小吧，倒也不是特别小，我觉得这事还是挺大的，不过不是特别大，但也不小。大家都觉得这事特别大，但我觉得也没那么大，但你要说小吧，这件事也不小。',
          ipdress: '上海',
          time: DateTime.now(),
        ),
        CommentPiece(
          avatar: 'https://s2.loli.net/2024/05/27/QzKM41C3Vs5FeHW.jpg',
          nickname: '风间澈',
          content: '我觉得火锅比较好吃耶，虽然火腿很香，有一点想吃mamamiya了哈哈哈，下次要一起去吗？',
          ipdress: '上海',
          time: DateTime.parse('2024-06-12 20:27:00'),
        ),
      ]),
    );
  }
}
