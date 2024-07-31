import 'package:flutter/material.dart';
import 'package:aorb/conf/config.dart';
import 'package:aorb/screens/follow_page.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/generated/user.pbgrpc.dart';

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
      body: Stack(
        children: [
          // 背景图片部分
          // TODO 添加用户个性化的设置，比如背景图片，消息卡片的颜色等
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            child: Container(
              height: MediaQuery.of(context).size.height * 0.4, // 设置背景图片的高度
              decoration: BoxDecoration(
                image: DecorationImage(
                  image: NetworkImage(user.avatar), // ~ 用户的背景图片设置
                  fit: BoxFit.cover,
                ),
              ),
            ),
          ),
          // 上层布局，包含左右布局
          Positioned.fill(
            child: Column(
              children: [
                Expanded(
                  flex: 3,
                  child: Row(
                    children: [
                      Expanded(
                        flex: 3,
                        child: Padding(
                          padding: const EdgeInsets.all(16.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Text(user.nickname, // ~ 用户昵称
                                  style: const TextStyle(
                                      color: Colors.white,
                                      fontSize: 32,
                                      fontWeight: FontWeight.bold,
                                      height: 1.5)),
                              Row(
                                children: <Widget>[
                                  Text(
                                    'Aorb ID: ',
                                    style: TextStyle(
                                      color: Colors.grey[100],
                                      fontSize: 12,
                                      height: 1.5,
                                    ),
                                  ),
                                  SelectableText(
                                    user.id, // ~ 用户ID
                                    style: TextStyle(
                                      color: Colors.grey[100],
                                      fontSize: 12,
                                      height: 1.5,
                                    ),
                                  ),
                                ],
                              ),
                              Text('IP归属地: ${user.ipaddress}', // ~ 用户IP归属地
                                  style: TextStyle(
                                      color: Colors.grey[100],
                                      fontSize: 12,
                                      height: 1.5)),
                              Row(
                                children: [
                                  GestureDetector(
                                    onTap: () {
                                      Navigator.push(
                                        context,
                                        MaterialPageRoute(
                                            builder: (context) =>
                                                FollowPage(username: user.id)),
                                      );
                                    },
                                    child: Text(
                                        '关注：${user.followed.usernames.length}', // ~ 用户关注的人数
                                        style: const TextStyle(
                                            color: Colors.white,
                                            fontSize: 12,
                                            height: 2)),
                                  ),
                                  const SizedBox(width: 16),
                                  GestureDetector(
                                    onTap: () {
                                      Navigator.push(
                                        context,
                                        MaterialPageRoute(
                                            builder: (context) =>
                                                FollowPage(username: user.id)),
                                      );
                                    },
                                    child: Text(
                                        '被关注：${user.follower.usernames.length}', // ~ 用户的粉丝数量
                                        style: const TextStyle(
                                            color: Colors.white,
                                            fontSize: 12,
                                            height: 2)),
                                  ),
                                ],
                              ),
                              Row(
                                children: [
                                  const Text('我的背包：',
                                      style: TextStyle(
                                          color: Colors.white,
                                          fontSize: 12,
                                          height: 2)),
                                  const Icon(
                                    // 金币图标
                                    Icons.monetization_on,
                                    color: Colors.white,
                                    size: 16,
                                  ),
                                  Text(user.coins.toString(), // ~ 用户的金币数量
                                      style: const TextStyle(
                                          color: Colors.white,
                                          fontSize: 12,
                                          height: 2)),
                                  const Spacer(),

                                  // TODO 调整“编辑资料”按钮的位置
                                  ElevatedButton(
                                    onPressed: () {}, // TODO 添加上编辑资料的逻辑
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor: Colors.blue[700],
                                      foregroundColor: Colors.white,
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(20),
                                      ),
                                    ),
                                    child: const Text('编辑资料'),
                                  ),
                                ],
                              ),
                            ],
                          ),
                        ),
                      ),
                      Expanded(
                        flex: 1,
                        child: Padding(
                          padding: const EdgeInsets.all(16.0),
                          child: LayoutBuilder(
                            builder: (context, constraints) {
                              double avatarSize =
                                  constraints.maxHeight / 6.5; // ^ 玄学调参
                              return CircleAvatar(
                                radius: avatarSize,
                                backgroundColor: Colors.white,
                                child: CircleAvatar(
                                  radius: avatarSize - 2, // 确保边框宽度
                                  backgroundImage:
                                      NetworkImage(user.avatar), // ~ 用户头像
                                  backgroundColor: Colors.transparent,
                                ),
                              );
                            },
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                // 回答问题部分
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
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                          unselectedLabelColor: Colors.grey[400],
                          indicator: UnderlineTabIndicator(
                              borderSide: BorderSide(
                                  width: 2.0, color: Colors.blue.shade700),
                              insets:
                                  const EdgeInsets.symmetric(horizontal: 20.0)),
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
                            // TODO 添加上从服务器拉取内容
                            children: const [
                              Center(child: Text('我发起的内容')),
                              Center(child: Text('我回答的内容')),
                              Center(child: Text('我收藏的内容')),
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
