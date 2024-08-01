import 'package:flutter/material.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/conf/config.dart';

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
  late User user; // 用户信息
  late List<User> followList; // 关注列表
  late List<User> fansList; // 粉丝列表

  @override
  void initState() {
    super.initState(); // 调用父类的 initState 方法
    _tabController = TabController(length: 2, vsync: this); // 初始化顶部导航栏控制器

    // 请求用户基本信息 关注列表 粉丝列表
    UserRequest request = UserRequest();
    request.username = widget.username;
    request.fields.addAll(['followed', 'follower']);
    UserService(backendHost, backendPort).getUserInfo(request).then((response) {
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

  @override
  Widget build(BuildContext context) {
    return Scaffold(
        appBar: AppBar(
          leading: IconButton(
            icon: const Icon(Icons.arrow_back),
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
              const Text('爱吃饭的小袁同学'),
            ],
          ),
        ),
        body: Column(
          children: [
            // 顶部导航栏
            TabBar(
              controller: _tabController,
              tabs: const [
                Tab(text: '关注'),
                Tab(text: '粉丝'),
              ],
            ),

            // 关注列表和粉丝列表
            // ~ 可能还需要修改
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  ListView.builder(
                    itemCount: followList.length,
                    itemBuilder: (context, index) {
                      return ListTile(
                        leading: CircleAvatar(
                          backgroundImage:
                              NetworkImage(followList[index].avatar),
                        ),
                        title: Text(followList[index].nickname),
                      );
                    },
                  ),
                  ListView.builder(
                    itemCount: fansList.length,
                    itemBuilder: (context, index) {
                      return ListTile(
                        leading: CircleAvatar(
                          backgroundImage: NetworkImage(fansList[index].avatar),
                        ),
                        title: Text(fansList[index].nickname),
                      );
                    },
                  ),
                ],
              ),
            ),
          ],
        ));
  }
}
