import 'package:flutter/material.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/services/user_service.dart';

// 页面查询并展示用户的关注列表和粉丝列表
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
    super.initState(); // 调用父类的 initState 方法
    _tabController = TabController(length: 2, vsync: this); // 初始化顶部导航栏控制器

    // 请求用户基本信息 关注列表 粉丝列表
    UserRequest request = UserRequest();
    request.username = widget.username;
    request.fields.addAll(['username','avatar', 'followed', 'follower']);
    UserService().getUserInfo(request).then((response) {
      setState(() {
        user = response.user;
        followList = user.followed;
        fansList = user.follower;
      });
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

// 根据username查询用户基本信息
  Future<User> getUserInfo(String username) async {
    UserRequest request = UserRequest();
    request.username = username;
    final response =
        await UserService().getUserInfo(request);
    return response.user;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        shadowColor: Colors.transparent, // 去掉 AppBar 阴影
        backgroundColor: Colors.white, // 修改 AppBar 底色为白色
        leading: IconButton(
          icon: const Icon(
            Icons.arrow_back,
            color: Colors.blue, // 修改返回按钮图标颜色为蓝色
          ),
          onPressed: () {
            Navigator.pop(context);
          },
        ),

        // 用户头像和用户名
        title: Row(
          children: [
            CircleAvatar(
              backgroundImage: NetworkImage(user.avatar),
              radius: 15,
            ),
            const SizedBox(width: 5),
            Text(
              widget.username, // 这里采用的是 follow_page 传入的 username
              style: const TextStyle(
                color: Colors.black,
                fontSize: 18,
              ), // 确保用户名颜色为黑色
            ),
          ],
        ),
      ),
      body: Column(
        children: [
          TabBar(
            controller: _tabController,
            // 标签样式
            labelColor: Colors.blue[700],
            labelStyle: const TextStyle(
              fontSize: 18,
              fontWeight: FontWeight.bold,
            ),
            unselectedLabelColor: Colors.grey[400],
            // 指示器样式
            indicator: UnderlineTabIndicator(
              borderSide: BorderSide(width: 2.0, color: Colors.blue.shade700),
              insets: const EdgeInsets.symmetric(horizontal: 20.0),
            ),
            indicatorSize: TabBarIndicatorSize.label,
            indicatorColor: Colors.blue[700],
            tabs: const [
              Tab(text: '关注'),
              Tab(text: '粉丝'),
            ],
          ),
          // 关注列表和粉丝列表
          Expanded(
            child: TabBarView(
              controller: _tabController,
              children: [
                // 关注列表
                ListView.builder(
                  itemCount: followList.usernames.length,
                  itemBuilder: (context, index) {
                    // 使用 FutureBuilder 处理异步获取用户信息
                    return FutureBuilder<User>(
                      future: getUserInfo(followList.usernames[index]),
                      builder: (context, snapshot) {
                        // 根据 snapshot 的状态显示不同的 UI
                        if (snapshot.connectionState ==
                            ConnectionState.waiting) {
                          // 数据正在加载时显示占位符
                          return const ListTile(
                            leading: CircleAvatar(
                              backgroundColor: Colors.grey, // 占位符颜色
                            ),
                            title: Text('Loading...'),
                          );
                        } else if (snapshot.hasError) {
                          // 出现错误时显示错误提示
                          return const ListTile(
                            leading: CircleAvatar(
                              backgroundColor: Colors.red, // 错误颜色
                            ),
                            title: Text('Error loading user info'),
                          );
                        } else if (snapshot.hasData) {
                          // 数据加载成功时显示用户信息
                          return ListTile(
                            leading: CircleAvatar(
                              backgroundImage:
                                  NetworkImage(snapshot.data!.avatar),
                            ),
                            title: Text(snapshot.data!.nickname),
                          );
                        } else {
                          // 没有数据时显示默认占位符
                          return const ListTile(
                            leading: CircleAvatar(
                              backgroundColor: Colors.grey, // 占位符颜色
                            ),
                            title: Text('No data available'),
                          );
                        }
                      },
                    );
                  },
                ),
                // 粉丝列表
                ListView.builder(
                  itemCount: fansList.usernames.length,
                  itemBuilder: (context, index) {
                    // 使用 FutureBuilder 处理异步获取用户信息
                    return FutureBuilder<User>(
                      future: getUserInfo(fansList.usernames[index]),
                      builder: (context, snapshot) {
                        // 根据 snapshot 的状态显示不同的 UI
                        if (snapshot.connectionState ==
                            ConnectionState.waiting) {
                          // 数据正在加载时显示占位符
                          return const ListTile(
                            leading: CircleAvatar(
                              backgroundColor: Colors.grey, // 占位符颜色
                            ),
                            title: Text('Loading...'),
                          );
                        } else if (snapshot.hasError) {
                          // 出现错误时显示错误提示
                          return const ListTile(
                            leading: CircleAvatar(
                              backgroundColor: Colors.red, // 错误颜色
                            ),
                            title: Text('Error loading user info'),
                          );
                        } else if (snapshot.hasData) {
                          // 数据加载成功时显示用户信息
                          return ListTile(
                            leading: CircleAvatar(
                              backgroundImage:
                                  NetworkImage(snapshot.data!.avatar),
                            ),
                            title: Text(snapshot.data!.nickname),
                          );
                        } else {
                          // 没有数据时显示默认占位符
                          return const ListTile(
                            leading: CircleAvatar(
                              backgroundColor: Colors.grey, // 占位符颜色
                            ),
                            title: Text('No data available'),
                          );
                        }
                      },
                    );
                  },
                ),
              ],
            ),
          )
        ],
      ),
    );
  }
}
