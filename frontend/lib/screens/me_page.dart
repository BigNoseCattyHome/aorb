import 'package:aorb/widgets/poll_list_view.dart';
import 'package:flutter/material.dart';

import 'package:aorb/screens/edit_profile_page.dart';
import 'package:aorb/screens/follow_page.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/services/user_service.dart';

import 'package:aorb/generated/user.pbgrpc.dart';
import 'package:aorb/generated/poll.pb.dart';

// 只有在登录的情况下才会展示这个页面，当没有登录的时候展示LoginPage，通过MainPage来检查登录状态并进行控制
class MePage extends StatefulWidget {
  final String username;
  const MePage({super.key, required this.username});

  @override
  MePageState createState() => MePageState();
}

class MePageState extends State<MePage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  User user = User();

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    // 根据 username 查询用户信息，并触发UI更新
    UserRequest request = UserRequest()..username = widget.username;
    UserService().getUserInfo(request).then((response) {
      setState(() {
        user = response.user;
      });
    });
  }

  void updateUserInfo(User updatedUser) {
    setState(() {
      user = updatedUser;
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  IconData _getGenderIcon(Gender gender) {
    switch (gender) {
      case Gender.MALE:
        return Icons.male;
      case Gender.FEMALE:
        return Icons.female;
      case Gender.OTHER:
        return Icons.transgender;
      default:
        return Icons.question_mark;
    }
  }

  Color _getGenderColor(Gender gender) {
    switch (gender) {
      case Gender.MALE:
        return Colors.blue;
      case Gender.FEMALE:
        return Colors.pink;
      case Gender.OTHER:
        return Colors.orange;
      default:
        return Colors.grey;
    }
  }

  Future<Map<String, dynamic>> _fetchPollData(String pollId) async {
    // 调用PollService获取投票信息
    final pollResponse = await PollService().GetPoll(
      GetPollRequest()..pollUuid = pollId,
    );
    final poll = pollResponse.poll;

    // 调用UserService获取用户信息
    final userInfoResponse = await UserService().getUserInfo(
      UserRequest()..username = poll.username,
    );
    final userInfo = userInfoResponse.user;

    // 计算总票数
    final totalVotes = poll.optionsCount.reduce((a, b) => a + b);

    // 计算每个选项的百分比
    List<double> percentages = poll.optionsCount.map((int value) {
      return (value / totalVotes) * 100;
    }).toList();

    return {
      'poll': poll,
      'userInfo': userInfo,
      'totalVotes': totalVotes,
      'percentages': percentages,
    };
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            child: Container(
              height: MediaQuery.of(context).size.height * 0.4,
              decoration: BoxDecoration(
                image: DecorationImage(
                  image: NetworkImage(user.bgpicMe),
                  fit: BoxFit.cover,
                  colorFilter: ColorFilter.mode(
                    Colors.black.withOpacity(0.4), // 透明度的黑色遮罩
                    BlendMode.darken,
                  ),
                ),
              ),
            ),
          ),
          SafeArea(
            child: Column(
              children: [
                // 个人信息部分
                Expanded(
                  flex: 3,
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Row(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // 左侧信息栏
                        Expanded(
                          flex: 3,
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              const SizedBox(height: 30),
                              // 用户昵称
                              Text(
                                user.nickname,
                                style: const TextStyle(
                                  color: Colors.white,
                                  fontSize: 32,
                                  fontWeight: FontWeight.bold,
                                  height: 1.5,
                                ),
                              ),
                              // Aorb ID
                              Row(
                                children: [
                                  Text(
                                    'Aorb ID: ',
                                    style: TextStyle(
                                      color: Colors.grey[100],
                                      fontSize: 12,
                                      height: 1.5,
                                    ),
                                  ),
                                  SelectableText(
                                    user.id,
                                    style: TextStyle(
                                      color: Colors.grey[100],
                                      fontSize: 12,
                                      height: 1.5,
                                    ),
                                  ),
                                ],
                              ),
                              // username
                              Text(
                                '用户名: ${user.username}',
                                style: TextStyle(
                                  color: Colors.grey[100],
                                  fontSize: 12,
                                  height: 1.5,
                                ),
                              ),

                              // Bio
                              Text(
                                'Bio: ${user.bio}',
                                style: TextStyle(
                                  color: Colors.grey[100],
                                  fontSize: 12,
                                  height: 1.5,
                                ),
                              ),

                              // IP归属地
                              Text(
                                'IP归属地: ${user.ipaddress}',
                                style: TextStyle(
                                  color: Colors.grey[100],
                                  fontSize: 12,
                                  height: 1.5,
                                ),
                              ),
                              const SizedBox(height: 25),
                              // 关注和被关注信息
                              Row(
                                children: [
                                  GestureDetector(
                                    onTap: () => Navigator.push(
                                      context,
                                      MaterialPageRoute(
                                        builder: (context) =>
                                            FollowPage(username: user.username),
                                      ),
                                    ),
                                    child: Text(
                                      '关注：${user.followed.usernames.length}',
                                      style: const TextStyle(
                                        color: Colors.white,
                                        fontSize: 12,
                                        height: 2,
                                      ),
                                    ),
                                  ),
                                  const SizedBox(width: 16),
                                  GestureDetector(
                                    onTap: () => Navigator.push(
                                      context,
                                      MaterialPageRoute(
                                        builder: (context) =>
                                            FollowPage(username: user.username),
                                      ),
                                    ),
                                    child: Text(
                                      '被关注：${user.follower.usernames.length}',
                                      style: const TextStyle(
                                        color: Colors.white,
                                        fontSize: 12,
                                        height: 2,
                                      ),
                                    ),
                                  ),
                                ],
                              ),
                              // 金币信息
                              Row(
                                children: [
                                  const Text(
                                    '我的背包：',
                                    style: TextStyle(
                                      color: Colors.white,
                                      fontSize: 12,
                                      height: 2,
                                    ),
                                  ),
                                  const Icon(
                                    Icons.monetization_on,
                                    color: Colors.white,
                                    size: 16,
                                  ),
                                  Text(
                                    user.coins.toString(),
                                    style: const TextStyle(
                                      color: Colors.white,
                                      fontSize: 12,
                                      height: 2,
                                    ),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                        // 右侧头像和编辑资料按钮
                        Expanded(
                          flex: 1,
                          child: Column(
                            children: [
                              const SizedBox(height: 50),
                              // 头像
                              CircleAvatar(
                                radius: 40,
                                backgroundColor: Colors.white,
                                child: Stack(
                                  children: [
                                    CircleAvatar(
                                      radius: 38,
                                      backgroundImage:
                                          NetworkImage(user.avatar),
                                      backgroundColor: Colors.transparent,
                                    ),
                                    Positioned(
                                      right: 0,
                                      bottom: 0,
                                      child: Container(
                                        padding: const EdgeInsets.all(2),
                                        decoration: const BoxDecoration(
                                          color: Colors.white,
                                          shape: BoxShape.circle,
                                        ),
                                        child: Icon(
                                          _getGenderIcon(user.gender),
                                          color: _getGenderColor(user.gender),
                                          size: 16,
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                              const SizedBox(height: 50),
                              // 编辑资料按钮
                              OutlinedButton(
                                onPressed: () async {
                                  final updatedUser = await Navigator.push(
                                    context,
                                    MaterialPageRoute(
                                      builder: (context) => EditProfilePage(
                                          user: user,
                                          onUserUpdated: updateUserInfo),
                                    ),
                                  );
                                  if (updatedUser != null) {
                                    updateUserInfo(updatedUser);
                                  }
                                },
                                style: OutlinedButton.styleFrom(
                                  foregroundColor: Colors.white,
                                  side: const BorderSide(color: Colors.white),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(20),
                                  ),
                                ),
                                child: const Text('编辑资料'),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
                // 下半部分（TabBar和TabBarView）
                Expanded(
                  flex: 5,
                  child: Container(
                    decoration: const BoxDecoration(
                      color: Colors.white,
                      borderRadius: BorderRadius.only(
                        topLeft: Radius.circular(20.0),
                        topRight: Radius.circular(20.0),
                      ),
                    ),
                    child: Column(
                      children: [
                        TabBar(
                          controller: _tabController,
                          labelColor: Colors.blue[700],
                          labelStyle: const TextStyle(
                            fontSize: 17,
                            fontWeight: FontWeight.bold,
                          ),
                          unselectedLabelColor: Colors.grey[400],
                          indicator: UnderlineTabIndicator(
                            borderSide: BorderSide(
                                width: 2.0, color: Colors.blue.shade700),
                            insets:
                                const EdgeInsets.symmetric(horizontal: 20.0),
                          ),
                          indicatorSize: TabBarIndicatorSize.label,
                          indicatorWeight: 3,
                          indicatorColor: Colors.blue[700],
                          tabs: const [
                            Tab(text: '我发起的'),
                            Tab(text: '我回答的'),
                            Tab(text: '我收藏的'),
                          ],
                        ),
                        Expanded(
                          child: TabBarView(
                            controller: _tabController,
                            children: [
                              PollListView(
                                pollIds: user.pollAsk.pollIds,
                                emptyMessage: '您还没有发起任何投票',
                              ),
                              PollListView(
                                pollIds: user.pollAns.pollIds,
                                emptyMessage: '您还没有回答任何投票',
                              ),
                              PollListView(
                                pollIds: user.pollCollect.pollIds,
                                emptyMessage: '您还没有收藏任何投票',
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
