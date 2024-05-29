import 'package:flutter/material.dart';

class MePage extends StatefulWidget {
  const MePage({super.key});

  @override
  MePageState createState() => MePageState();
}

class MePageState extends State<MePage> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
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
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            child: Container(
              height: MediaQuery.of(context).size.height * 0.4, // 设置背景图片的高度
              decoration: const BoxDecoration(
                image: DecorationImage(
                  image: NetworkImage(
                      'https://s2.loli.net/2024/05/28/lTEVmXLvwuQW9in.jpg'),
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
                              const Text('大鼻猫之家',
                                  style: TextStyle(
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
                                    'bignosecattyhome',
                                    style: TextStyle(
                                      color: Colors.grey[100],
                                      fontSize: 12,
                                      height: 1.5,
                                    ),
                                  ),
                                ],
                              ),
                              Text('IP归属地: 上海',
                                  style: TextStyle(
                                      color: Colors.grey[100],
                                      fontSize: 12,
                                      height: 1.5)),
                              const Row(
                                children: [
                                  Text('关注：2',
                                      style: TextStyle(
                                          color: Colors.white,
                                          fontSize: 12,
                                          height: 2)),
                                  SizedBox(width: 16),
                                  Text('被关注：2',
                                      style: TextStyle(
                                          color: Colors.white,
                                          fontSize: 12,
                                          height: 2)),
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
                                  const Text('30',
                                      style: TextStyle(
                                          color: Colors.white,
                                          fontSize: 12,
                                          height: 2)),
                                  const Spacer(),
                                  ElevatedButton(
                                    onPressed: () {},
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
                                  backgroundImage: const NetworkImage(
                                      'https://s2.loli.net/2024/05/28/MY6bk5FxVh8ufOa.png'),
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
                            borderSide: BorderSide(width: 2.0, color: Colors.blue.shade700),
                            insets: const EdgeInsets.symmetric(horizontal: 20.0)
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
