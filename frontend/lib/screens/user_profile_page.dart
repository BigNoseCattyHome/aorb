import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/generated/user.pbgrpc.dart';
import 'package:flutter/material.dart';
import 'package:aorb/widgets/poll_list_view.dart';
import 'package:aorb/screens/follow_page.dart';
import 'package:aorb/services/user_service.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:shared_preferences/shared_preferences.dart';

class UserProfilePage extends StatefulWidget {
  final String username;

  const UserProfilePage({super.key, required this.username});

  @override
  UserProfilePageState createState() => UserProfilePageState();
}

class UserProfilePageState extends State<UserProfilePage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  User user = User();
  final logger = getLogger();
  String currentUsername = "";

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _fetchUserInfo();
    // 从shared_preferences中获取当前用户的用户名
    SharedPreferences.getInstance().then((prefs) {
      currentUsername = prefs.getString('username') ?? '';
    });
  }

  void _fetchUserInfo() {
    UserRequest request = UserRequest()
      ..username = widget.username
      ..fields.addAll(
          ["avatar", "bgpicMe", "nickname", "username", "bio", "ipaddress"]);
    UserService().getUserInfo(request).then((response) {
      setState(() {
        user = response.user;
      });
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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          _buildBackgroundImage(),
          SafeArea(
            child: Column(
              children: [
                Expanded(
                  flex: 3,
                  child: _buildUserInfoSection(),
                ),
                Expanded(
                  flex: 5,
                  child: _buildTabSection(),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBackgroundImage() {
    return Positioned(
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
              Colors.black.withOpacity(0.4),
              BlendMode.darken,
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildUserInfoSection() {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Expanded(
            flex: 3,
            child: _buildLeftInfoColumn(),
          ),
          Expanded(
            flex: 1,
            child: _buildRightAvatarColumn(),
          ),
        ],
      ),
    );
  }

  Widget _buildLeftInfoColumn() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: 30),
        _buildNickname(),
        _buildUsername(),
        _buildBio(),
        _buildIpAddress(),
        const SizedBox(height: 25),
        _buildFollowInfo(),
      ],
    );
  }

  Widget _buildNickname() {
    return Text(
      user.nickname,
      style: const TextStyle(
        color: Colors.white,
        fontSize: 32,
        fontWeight: FontWeight.bold,
        height: 1.5,
      ),
    );
  }

  Widget _buildUsername() {
    return Text(
      '用户名: ${user.username}',
      style: TextStyle(
        color: Colors.grey[100],
        fontSize: 12,
        height: 1.5,
      ),
    );
  }

  Widget _buildBio() {
    return Text(
      'Bio: ${user.bio}',
      style: TextStyle(
        color: Colors.grey[100],
        fontSize: 12,
        height: 1.5,
      ),
    );
  }

  Widget _buildIpAddress() {
    return Text(
      'IP归属地: ${user.ipaddress}',
      style: TextStyle(
        color: Colors.grey[100],
        fontSize: 12,
        height: 1.5,
      ),
    );
  }

  Widget _buildFollowInfo() {
    return Row(
      children: [
        GestureDetector(
          onTap: () => _navigateToFollowPage(),
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
          onTap: () => _navigateToFollowPage(),
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
    );
  }

  void _navigateToFollowPage() {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => FollowPage(username: user.username),
      ),
    );
  }

  Widget _buildRightAvatarColumn() {
    return Column(
      children: [
        const SizedBox(height: 50),
        _buildAvatar(),
        const SizedBox(height: 20),
        _buildFollowButton(),
      ],
    );
  }

  Widget _buildFollowButton() {
    bool isFollowing = user.follower.usernames.contains(currentUsername);

    return ElevatedButton(
      onPressed: _toggleFollow,
      style: ElevatedButton.styleFrom(
        foregroundColor: isFollowing ? Colors.grey[800] : Colors.white,
        backgroundColor: isFollowing ? Colors.grey[300] : Colors.blue,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
        ),
        padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 10),
      ),
      child: Text(
        isFollowing ? '取消关注' : '关注',
        style: const TextStyle(fontSize: 14),
      ),
    );
  }

  void _toggleFollow() async {
    setState(() {
      // 立即更新UI以提供即时反馈
      if (user.follower.usernames.contains(currentUsername)) {
        user.follower.usernames.remove(currentUsername);
      } else {
        user.follower.usernames.add(currentUsername);
      }
    });

    try {
      // 发送关注/取消关注请求到服务器
      final request = FollowUserRequest()
        ..username = currentUsername
        ..targetUsername = user.username;

      final response = await UserService().followUser(request);

      if (response.statusCode == 0) {
        // 如果服务器操作成功,不需要做任何事,因为我们已经更新了UI
      } else {
        // 如果服务器操作失败,恢复原来的状态
        setState(() {
          if (user.follower.usernames.contains(currentUsername)) {
            user.follower.usernames.remove(currentUsername);
          } else {
            user.follower.usernames.add(currentUsername);
          }
        });

        // 显示错误消息
        Fluttertoast.showToast(msg: '操作失败: ${response.statusMsg}');
      }
    } catch (e) {
      // 处理网络错误等异常情况
      logger.e('Error toggling follow: $e');
      // 恢复原来的状态
      setState(() {
        if (user.follower.usernames.contains(currentUsername)) {
          user.follower.usernames.remove(currentUsername);
        } else {
          user.follower.usernames.add(currentUsername);
        }
      });
      // 显示错误消息
      Fluttertoast.showToast(msg: '网络错误,请稍后再试');
    }
  }

  Widget _buildAvatar() {
    return CircleAvatar(
      radius: 40,
      backgroundColor: Colors.white,
      child: Stack(
        children: [
          CircleAvatar(
            radius: 38,
            backgroundImage: NetworkImage(user.avatar),
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
    );
  }

  Widget _buildTabSection() {
    return Container(
      decoration: const BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.only(
          topLeft: Radius.circular(20.0),
          topRight: Radius.circular(20.0),
        ),
      ),
      child: Column(
        children: [
          _buildTabBar(),
          Expanded(
            child: _buildTabBarView(),
          ),
        ],
      ),
    );
  }

  Widget _buildTabBar() {
    return TabBar(
      controller: _tabController,
      labelColor: Colors.blue[700],
      labelStyle: const TextStyle(
        fontSize: 17,
        fontWeight: FontWeight.bold,
      ),
      unselectedLabelColor: Colors.grey[400],
      indicator: UnderlineTabIndicator(
        borderSide: BorderSide(width: 2.0, color: Colors.blue.shade700),
        insets: const EdgeInsets.symmetric(horizontal: 20.0),
      ),
      indicatorSize: TabBarIndicatorSize.label,
      indicatorWeight: 3,
      indicatorColor: Colors.blue[700],
      tabs: const [
        Tab(text: '发起的投票'),
        Tab(text: '回答的投票'),
        Tab(text: '收藏的投票'),
      ],
    );
  }

  Widget _buildTabBarView() {
    return TabBarView(
      controller: _tabController,
      children: [
        PollListView(
          pollIds: user.pollAsk.pollIds,
          emptyMessage: '该用户还没有发起任何投票',
          currentUsername: currentUsername,
        ),
        PollListView(
          pollIds: user.pollAns.pollIds,
          emptyMessage: '该用户还没有回答任何投票',
          currentUsername: currentUsername,
        ),
        PollListView(
          pollIds: user.pollCollect.pollIds,
          emptyMessage: '该用户还没有收藏任何投票',
          currentUsername: currentUsername,
        ),
      ],
    );
  }
}
