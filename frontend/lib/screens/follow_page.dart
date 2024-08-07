import 'package:flutter/material.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/services/user_service.dart';

// 关注页面
// 展示用户的关注列表和粉丝列表
class FollowPage extends StatefulWidget {
  final String username;
  const FollowPage({Key? key, required this.username}) : super(key: key);

  @override
  FollowPageState createState() => FollowPageState();
}

class FollowPageState extends State<FollowPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController; // 顶部导航栏控制器
  User user = User(); // 用户信息
  FollowedList followList = FollowedList(); // 关注列表
  FollowerList fansList = FollowerList(); // 粉丝列表

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
    _fetchUserInfo();
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  // 获取用户信息
  void _fetchUserInfo() {
    UserRequest request = UserRequest()
      ..username = widget.username
      ..fields.addAll(['username', 'avatar', 'followed', 'follower']);

    UserService().getUserInfo(request).then((response) {
      setState(() {
        user = response.user;
        followList = user.followed;
        fansList = user.follower;
      });
    });
  }

  // 根据username查询用户基本信息
  Future<User> getUserInfo(String username) async {
    UserRequest request = UserRequest()..username = username;
    final response = await UserService().getUserInfo(request);
    return response.user;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: _buildAppBar(),
      body: Column(
        children: [
          _buildTabBar(),
          Expanded(
            child: _buildTabBarView(),
          ),
        ],
      ),
    );
  }

  // 构建AppBar
  AppBar _buildAppBar() {
    return AppBar(
      shadowColor: Colors.transparent,
      backgroundColor: Colors.white,
      leading: IconButton(
        icon: const Icon(Icons.arrow_back, color: Colors.blue),
        onPressed: () => Navigator.pop(context),
      ),
      title: Row(
        children: [
          CircleAvatar(
            backgroundImage: NetworkImage(user.avatar),
            radius: 15,
          ),
          const SizedBox(width: 5),
          Text(
            widget.username,
            style: const TextStyle(color: Colors.black, fontSize: 18),
          ),
        ],
      ),
    );
  }

  // 构建TabBar
  Widget _buildTabBar() {
    return TabBar(
      controller: _tabController,
      labelColor: Colors.blue[700],
      labelStyle: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
      unselectedLabelColor: Colors.grey[400],
      indicator: UnderlineTabIndicator(
        borderSide: BorderSide(width: 2.0, color: Colors.blue.shade700),
        insets: const EdgeInsets.symmetric(horizontal: 20.0),
      ),
      indicatorSize: TabBarIndicatorSize.label,
      indicatorColor: Colors.blue[700],
      tabs: const [Tab(text: '关注'), Tab(text: '粉丝')],
    );
  }

  // 构建TabBarView
  Widget _buildTabBarView() {
    return TabBarView(
      controller: _tabController,
      children: [
        _buildUserList(followList.usernames),
        _buildUserList(fansList.usernames),
      ],
    );
  }

  // 构建用户列表
  Widget _buildUserList(List<String> usernames) {
    return ListView.builder(
      itemCount: usernames.length,
      itemBuilder: (context, index) {
        return _buildUserListItem(usernames[index]);
      },
    );
  }

  // 构建用户列表项
  Widget _buildUserListItem(String username) {
    return FutureBuilder<User>(
      future: getUserInfo(username),
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return _buildLoadingListTile();
        } else if (snapshot.hasError) {
          return _buildErrorListTile();
        } else if (snapshot.hasData) {
          return _buildUserInfoListTile(snapshot.data!);
        } else {
          return _buildNoDataListTile();
        }
      },
    );
  }

  // 构建加载中的ListTile
  Widget _buildLoadingListTile() {
    return const ListTile(
      leading: CircleAvatar(backgroundColor: Colors.grey),
      title: Text('Loading...'),
    );
  }

  // 构建错误的ListTile
  Widget _buildErrorListTile() {
    return const ListTile(
      leading: CircleAvatar(backgroundColor: Colors.red),
      title: Text('Error loading user info'),
    );
  }

  // 构建用户信息的ListTile
  Widget _buildUserInfoListTile(User user) {
    return ListTile(
      leading: CircleAvatar(backgroundImage: NetworkImage(user.avatar)),
      title: Text(user.nickname),
    );
  }

  // 构建无数据的ListTile
  Widget _buildNoDataListTile() {
    return const ListTile(
      leading: CircleAvatar(backgroundColor: Colors.grey),
      title: Text('No data available'),
    );
  }
}