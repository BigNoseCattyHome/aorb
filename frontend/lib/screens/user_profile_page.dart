import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/generated/user.pbgrpc.dart';
import 'package:flutter/material.dart';
import 'package:aorb/widgets/poll_list_view.dart';
import 'package:aorb/screens/follow_page.dart';
import 'package:aorb/services/user_service.dart';
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
  bool isFollowed = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    _fetchUserInfo();
    SharedPreferences.getInstance().then((prefs) {
      setState(() {
        currentUsername = prefs.getString('username') ?? '';
      });
      _checkFollowStatus();
    });
  }

  void _fetchUserInfo() {
    UserRequest request = UserRequest()..username = widget.username;
    UserService().getUserInfo(request).then((response) {
      setState(() {
        user = response.user;
      });
    });
  }

  Future<void> _checkFollowStatus() async {
    if (currentUsername.isEmpty) {
      // 如果 currentUsername 还没有加载，等待它加载
      final prefs = await SharedPreferences.getInstance();
      currentUsername = prefs.getString('username') ?? '';
    }
    var request = IsUserFollowingRequest()
      ..username = currentUsername
      ..targetUsername = widget.username;
    final isFollowing = await UserService().isUserFollowing(request);
    setState(() {
      isFollowed = isFollowing;
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
    return Visibility(
      visible: user.username != currentUsername,
      child: TextButton(
        onPressed: _toggleFollow,
        style: TextButton.styleFrom(
          backgroundColor: isFollowed ? Colors.white : Colors.red,
          foregroundColor: isFollowed ? Colors.blue : Colors.white,
          side:
              BorderSide(color: isFollowed ? Colors.grey : Colors.transparent),
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(20),
          ),
          minimumSize: const Size(80, 30),
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
        ),
        child: Text(
          isFollowed ? '已关注' : '关注',
          style: const TextStyle(
            fontSize: 12,
          ),
        ),
      ),
    );
  }

  Future<void> _toggleFollow() async {
    try {
      if (isFollowed) {
        var request = FollowUserRequest()
          ..username = currentUsername
          ..targetUsername = widget.username;
        await UserService().unfollowUser(request);
      } else {
        var request = FollowUserRequest()
          ..username = currentUsername
          ..targetUsername = user.username;
        await UserService().followUser(request);
      }
      setState(() {
        isFollowed = !isFollowed;
        if (isFollowed) {
          user.follower.usernames.add(currentUsername);
        } else {
          user.follower.usernames.remove(currentUsername);
        }
      });
    } catch (e) {
      logger.e('Failed to follow: $e');
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

  void refreshPollListViews() {
    setState(() {
      // 这里不需要做任何事情，因为setState会触发重建
    });
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
          onRefresh: refreshPollListViews,
        ),
        PollListView(
          pollIds: user.pollAns.pollIds,
          emptyMessage: '该用户还没有回答任何投票',
          currentUsername: currentUsername,
          onRefresh: refreshPollListViews,
        ),
        PollListView(
          pollIds: user.pollCollect.pollIds,
          emptyMessage: '该用户还没有收藏任何投票',
          currentUsername: currentUsername,
          onRefresh: refreshPollListViews,
        ),
      ],
    );
  }
}
