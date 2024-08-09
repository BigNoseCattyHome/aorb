import 'package:flutter/material.dart';
import 'package:aorb/generated/poll.pb.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/widgets/say_something.dart';
import 'package:aorb/widgets/comment_piece.dart';
import 'package:aorb/widgets/poll_detail.dart';
import 'package:aorb/generated/user.pbgrpc.dart';

class PollDetailPage extends StatefulWidget {
  final String userId;
  final String pollId;
  final String username;
  const PollDetailPage(
      {super.key,
      required this.pollId,
      required this.userId,
      required this.username});

  @override
  PollDetailPageState createState() => PollDetailPageState();
}

class PollDetailPageState extends State<PollDetailPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  bool isFollowed = false;
  int cntComments = 0;
  Poll? poll;
  User? user;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadData();
  }

  Future<void> _loadData() async {
    // 加载关注状态
    final requestFollow = IsUserFollowingRequest()..username = widget.username;
    final isFollowed = await UserService().isUserFollowing(requestFollow);

    // 加载用户信息
    final userRequest = UserRequest()..username = widget.username;
    final userResponse = await UserService().getUserInfo(userRequest);

    // 加载投票详情
    final pollResponse =
        await PollService().getPoll(GetPollRequest()..pollUuid = widget.pollId);

    setState(() {
      this.isFollowed = isFollowed;
      user = userResponse.user;
      poll = pollResponse.poll;
      cntComments = poll?.commentList.length ?? 0;
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.pop(context),
        ),
        title: Row(
          children: [
            if (user != null) ...[
              CircleAvatar(
                backgroundImage: NetworkImage(user!.avatar),
                radius: 15,
              ),
              const SizedBox(width: 5),
              Text(user!.nickname),
            ],
          ],
        ),
        backgroundColor: Colors.white,
        elevation: 0,
        actions: _buildFollowButton(),
      ),
      body: poll == null
          ? const Center(child: CircularProgressIndicator())
          : SingleChildScrollView(
              child: Column(
                children: [
                  // 问题详情
                  PollDetail(
                    title: poll!.title,
                    content: poll!.content,
                    options: poll!.options,
                    votePercentage: poll!.optionsCount.map((value) {
                      final total = poll!.optionsCount.reduce((a, b) => a + b);
                      return total > 0 ? (value / total) * 100 : 0.0;
                    }).toList(),
                    time: poll!.createAt.toDateTime(),
                    ipaddress: poll!.ipaddress,
                  ),

                  // 分割线
                  const Divider(),
                  Text('共有$cntComments条评论'),
                  const SaySomething(avatar: 'avatar'),

                  // 评论列表
                  FutureBuilder<List<Widget>>(
                    future: Future.wait(poll!.commentList.map((comment) async {
                      // 获取评论者的用户信息
                      final userRequest = UserRequest()
                        ..username = comment.commentUsername;
                      final userResponse =
                          await UserService().getUserInfo(userRequest);
                      final user = userResponse.user;

                      return CommentPiece(
                        content: comment.content,
                        time: comment.createAt.toDateTime(),
                        avatar: user.avatar,
                        nickname: user.nickname,
                        ipdress: user.ipaddress,
                      );
                    })),
                    builder: (context, snapshot) {
                      if (snapshot.connectionState == ConnectionState.waiting) {
                        return const Center(child: CircularProgressIndicator());
                      } else if (snapshot.hasError) {
                        return Center(child: Text('Error: ${snapshot.error}'));
                      } else if (snapshot.hasData) {
                        return ListView(
                          shrinkWrap: true,
                          physics: const NeverScrollableScrollPhysics(),
                          children: snapshot.data!,
                        );
                      } else {
                        return const Center(child: Text('No comments'));
                      }
                    },
                  ),
                ],
              ),
            ),
    );
  }

  List<Widget> _buildFollowButton() {
    return [
      TextButton(
        style: TextButton.styleFrom(
          side: const BorderSide(color: Colors.blue),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
        ),
        onPressed: () {
          // TODO: 实现关注/取消关注逻辑
          setState(() {
            isFollowed = !isFollowed;
          });
        },
        child: Text(
          isFollowed ? '已关注' : '关注',
          style: const TextStyle(color: Colors.blue),
        ),
      ),
    ];
  }
}
