import 'package:aorb/models/user.dart';
import 'package:flutter/material.dart';
import 'package:aorb/models/simple_user.dart';
import 'package:aorb/services/user_service.dart';


// 页面展示用户ID为userID的关注列表和粉丝列表
class FollowPage extends StatefulWidget {
  final String userId;
  const FollowPage({Key? key, required this.userId}) : super(key: key);

  @override
  FollowPageState createState() => FollowPageState();
}

class FollowPageState extends State<FollowPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController; // 顶部导航栏控制器
  late User user; // 用户信息
  late List<SimpleUser> followList; // 关注列表
  late List<SimpleUser> fansList; // 粉丝列表

  @override
  void initState() {
    super.initState(); // 调用父类的 initState 方法
    _tabController = TabController(length: 2, vsync: this); // 初始化顶部导航栏控制器

    // 请求用户nickname和avtar
    UserService()
        .fetchUserInfo(widget.userId, ['nickname', 'avtar']).then((value) {
      setState(() {
        user = value;
      });
    });

    // 获取关注列表
    UserService().fetchFollowList(widget.userId).then((value) {
      setState(() {
        followList = value;
      });
    });

    // 获取粉丝列表
    UserService().fetchFanList(widget.userId).then((value) {
      setState(() {
        fansList = value;
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
                        backgroundImage: NetworkImage(followList[index].avatar),
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
      )
    );
  }
}
