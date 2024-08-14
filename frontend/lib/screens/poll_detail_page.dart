import 'package:aorb/conf/config.dart';
import 'package:aorb/screens/user_profile_page.dart';
import 'package:flutter/material.dart';
import 'package:aorb/generated/poll.pb.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/widgets/say_something.dart';
import 'package:aorb/widgets/comment_piece.dart';
import 'package:aorb/widgets/poll_detail.dart';
import 'package:aorb/generated/user.pbgrpc.dart';
import 'package:shared_preferences/shared_preferences.dart';

class PollDetailPage extends StatefulWidget {
  final String userId; // 暂时不用
  final String pollId;
  final String username; // 问题发布者的用户名

  const PollDetailPage({
    super.key,
    required this.pollId,
    required this.userId,
    required this.username,
  });

  @override
  PollDetailPageState createState() => PollDetailPageState();
}

class PollDetailPageState extends State<PollDetailPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  bool isFollowed = false;
  int cntComments = 0;
  Poll poll = Poll();
  User user = User();
  String currentUser = '';
  String currentUserAvatar = '';
  late SharedPreferences prefs;
  String selectedOption = "";
  final logger = getLogger();

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _loadData();
  }

  // 加载数据的方法
  Future<void> _loadData() async {
    // 获取当前用户名，从本地存储中获取
    prefs = await SharedPreferences.getInstance();
    currentUser = prefs.getString('username') ?? '';

    // 加载关注状态
    final requestFollow = IsUserFollowingRequest()
      ..username = currentUser
      ..targetUsername = widget.username;
    isFollowed = await UserService().isUserFollowing(requestFollow);

    // 加载用户信息
    final userRequest = UserRequest()..username = widget.username;
    final userResponse = await UserService().getUserInfo(userRequest);
    final currentUserRequest = UserRequest()
      ..username = currentUser
      ..fields.addAll(['username', 'avatar']);
    final currentUserResponse =
        await UserService().getUserInfo(currentUserRequest);

    // 加载投票详情
    final pollResponse =
        await PollService().getPoll(GetPollRequest()..pollUuid = widget.pollId);

    // 查询用户是否已经投票
    final selectedOptionResponse =
        await PollService().getChoiceWithPollUuidAndUsername(
      GetChoiceWithPollUuidAndUsernameRequest()
        ..pollUuid = widget.pollId
        ..username = widget.username,
    );

    if (mounted) {
      setState(() {
        isFollowed = isFollowed;
        user = userResponse.user;
        poll = pollResponse.poll;
        cntComments = poll.commentList.length;
        selectedOption = selectedOptionResponse.choice;
        currentUserAvatar = currentUserResponse.user.avatar;
      });
    } else {
      logger.e('Error loading poll detail when PollDetailPage is not mounted');
    }
  }

  // 刷新数据的方法
  Future<void> _refreshData() async {
    await _loadData();
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
          color: Colors.black,
          onPressed: () => Navigator.pop(context),
        ),
        title: GestureDetector(
          onTap: () {
            // 跳转到用户详情页面
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => UserProfilePage(username: user.username),
              ),
            );
          },
          child: Row(
            children: [
              CircleAvatar(
                backgroundImage: NetworkImage(user.avatar),
                radius: 15,
              ),
              const SizedBox(width: 5),
              Text(
                user.nickname,
                style: const TextStyle(color: Colors.black),
              ),
            ],
          ),
        ),
        backgroundColor: Colors.white,
        elevation: 0,
        actions: [
          Center(
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 8.0),
              child: _buildFollowButton(),
            ),
          ),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _refreshData,
        child: poll == Poll()
            ? const Center(child: CircularProgressIndicator())
            : SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                child: Column(
                  children: [
                    // 问题详情
                    PollDetail(
                      title: poll.title,
                      content: poll.content,
                      options: poll.options,
                      votePercentage: poll.optionsCount.map((value) {
                        final total = poll.optionsCount.reduce((a, b) => a + b);
                        return total > 0 ? (value / total) * 100 : 0.0;
                      }).toList(),
                      time: poll.createAt,
                      ipaddress: poll.ipaddress,
                      selectedOption: selectedOption,
                      pollId: poll.pollUuid,
                      username: user.username,
                      bgpic: user.bgpicPollcard,
                      currentUser: currentUser,
                    ),

                    // 分割线
                    const Divider(),

                    // 评论区
                    SaySomething(
                      currentUserAvatar: currentUserAvatar,
                      currentUsername: currentUser,
                      pollId: poll.pollUuid,
                      onCommentPosted: _refreshData, // 添加评论后刷新
                    ),
                    const SizedBox(height: 10),
                    cntComments > 0
                        ? Text(
                            '共有$cntComments条评论',
                            style: const TextStyle(color: Colors.grey),
                            textAlign: TextAlign.left,
                          )
                        : const Text(
                            '评论区还是一片荒草地，快来抢占沙发吧！',
                            style: TextStyle(color: Colors.grey),
                            textAlign: TextAlign.left,
                          ),

                    // 评论列表
                    FutureBuilder<List<Widget>>(
                      future: Future.wait(poll.commentList.map((comment) async {
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
                        if (snapshot.connectionState ==
                            ConnectionState.waiting) {
                          return const Center(
                              child: CircularProgressIndicator());
                        } else if (snapshot.hasError) {
                          return Center(
                              child: Text('Error: ${snapshot.error}'));
                        } else if (snapshot.hasData) {
                          return ListView(
                            shrinkWrap: true,
                            physics: const NeverScrollableScrollPhysics(),
                            children: snapshot.data!,
                          );
                        } else {
                          return const Center(child: Text('暂无评论'));
                        }
                      },
                    ),
                  ],
                ),
              ),
      ),
    );
  }

  Widget _buildFollowButton() {
    return Visibility(
      visible: user.username != currentUser,
      child: TextButton(
        style: TextButton.styleFrom(
          backgroundColor: isFollowed ? Colors.transparent : Colors.red,
          side:
              BorderSide(color: isFollowed ? Colors.grey : Colors.transparent),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          minimumSize: const Size(60, 30),
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
        ),
        onPressed: () {
          if (isFollowed) {
            var request = FollowUserRequest()
              ..username = currentUser
              ..targetUsername = widget.username;
            UserService().unfollowUser(request);
          } else {
            var request = FollowUserRequest()
              ..username = currentUser
              ..targetUsername = user.username;
            UserService().followUser(request);
          }
          setState(() {
            isFollowed = !isFollowed;
          });
        },
        child: Text(
          isFollowed ? '已关注' : '关注',
          style: TextStyle(
            color: isFollowed ? Colors.grey : Colors.white,
            fontSize: 12,
          ),
        ),
      ),
    );
  }
}
