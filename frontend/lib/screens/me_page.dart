import 'package:aorb/widgets/poll_list_view.dart';
import 'package:flutter/material.dart';
import 'package:aorb/screens/edit_profile_page.dart';
import 'package:aorb/screens/follow_page.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/generated/user.pbgrpc.dart';

class MePage extends StatefulWidget {
  final String username;
  final Function(String) onAvatarUpdated; // 新增：用于通知 MainPage 头像已更新

  const MePage(
      {super.key, required this.username, required this.onAvatarUpdated});

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
    _fetchUserInfo();
  }

  // 获取用户信息
  void _fetchUserInfo() {
    UserRequest request = UserRequest()..username = widget.username;
    UserService().getUserInfo(request).then((response) {
      setState(() {
        user = response.user;
      });
    });
  }

  // 更新用户信息
  void updateUserInfo(User updatedUser) {
    setState(() {
      user = updatedUser;
    });
    // 如果头像更新了，通知 MainPage
    if (user.avatar != widget.onAvatarUpdated.toString()) {
      widget.onAvatarUpdated(user.avatar);
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  // 获取性别图标
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

  // 获取性别颜色
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

  Future<void> _onRefresh() async {
    _fetchUserInfo();
    setState(() {});
  }

  void refreshPollListViews() {
    setState(() {});
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: RefreshIndicator(
        onRefresh: _onRefresh,
        child: Stack(
          children: [
            _buildBackgroundImage(),
            SafeArea(
              child: SingleChildScrollView(
                physics: const AlwaysScrollableScrollPhysics(),
                child: SizedBox(
                  height: MediaQuery.of(context).size.height,
                  child: Column(
                    children: [
                      Expanded(
                        flex: 3,
                        child: _buildUserInfoSection(),
                      ),
                      Expanded(
                        flex: 6,
                        child: _buildTabSection(),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  // 构建背景图像
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

  // 构建用户信息部分
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

  // 构建左侧信息栏
  Widget _buildLeftInfoColumn() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SizedBox(height: 30),
        _buildNickname(),
        _buildAorbId(),
        _buildUsername(),
        _buildBio(),
        _buildIpAddress(),
        const SizedBox(height: 25),
        _buildFollowInfo(),
        _buildCoinsInfo(),
      ],
    );
  }

  // 构建昵称
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

  // 构建Aorb ID
  Widget _buildAorbId() {
    return Row(
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
    );
  }

  // 构建用户名
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

  // 构建个人简介
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

  // 构建IP地址
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

  // 构建关注信息
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

  // 导航到关注页面
  void _navigateToFollowPage() {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => FollowPage(username: user.username),
      ),
    );
  }

  // 构建金币信息
  Widget _buildCoinsInfo() {
    return Row(
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
    );
  }

  // 构建右侧头像栏
  Widget _buildRightAvatarColumn() {
    return Column(
      children: [
        const SizedBox(height: 50),
        _buildAvatar(),
        const SizedBox(height: 50),
        _buildEditProfileButton(),
      ],
    );
  }

  // 构建头像
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

  // 构建编辑资料按钮
  Widget _buildEditProfileButton() {
    return OutlinedButton(
      onPressed: _navigateToEditProfilePage,
      style: OutlinedButton.styleFrom(
        foregroundColor: Colors.white,
        side: const BorderSide(color: Colors.white),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(20),
        ),
      ),
      child: const Text('编辑资料'),
    );
  }

  // 导航到编辑资料页面
  void _navigateToEditProfilePage() async {
    final updatedUser = await Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => EditProfilePage(
          user: user,
          onUserUpdated: updateUserInfo,
        ),
      ),
    );
    if (updatedUser != null) {
      updateUserInfo(updatedUser);
    }
  }

  // 构建标签页部分
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

  // 构建标签栏
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
        Tab(text: '我发起的'),
        Tab(text: '我回答的'),
        Tab(text: '我收藏的'),
      ],
    );
  }

  // 构建标签页视图
  Widget _buildTabBarView() {
    return TabBarView(
      controller: _tabController,
      children: [
        PollListView(
          key: UniqueKey(),
          pollIds: user.pollAsk.pollIds,
          emptyMessage: '您还没有发起任何投票',
          currentUsername: widget.username,
          onRefresh: refreshPollListViews,
        ),
        PollListView(
          key: UniqueKey(),
          pollIds: user.pollAns.pollIds,
          emptyMessage: '您还没有回答任何投票',
          currentUsername: widget.username,
          onRefresh: refreshPollListViews,
        ),
        PollListView(
          key: UniqueKey(),
          pollIds: user.pollCollect.pollIds,
          emptyMessage: '您还没有收藏任何投票',
          currentUsername: widget.username,
          onRefresh: refreshPollListViews,
        ),
      ],
    );
  }
}
