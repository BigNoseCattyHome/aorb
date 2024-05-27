import 'package:flutter/material.dart';

class DynamicTopBar extends StatelessWidget implements PreferredSizeWidget {
  final List<String> tabs;
  final bool showSearch;
  final TabController tabController; // 接受外部提供的TabController

  const DynamicTopBar({
    Key? key,
    required this.tabs,
    this.showSearch = true,
    required this.tabController, // 现在是必需的属性
  }) : super(key: key);

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      elevation: 0,
      backgroundColor: Colors.white,
      leading: IconButton(
        icon: Icon(Icons.menu, color: Colors.blue[700]),
        onPressed: () {
          Scaffold.of(context).openDrawer();
        },
      ),
      title: Row(
        children: [
          const SizedBox(width: 25),
          Expanded(
            child: TabBar(
              controller: tabController, // 使用传入的TabController
              tabs: tabs.map((tab) => Tab(text: tab)).toList(),
              labelColor: Colors.blue[700],
              labelStyle: const TextStyle(
                fontSize: 20,
                fontFamily: 'SimHei',
                fontWeight: FontWeight.bold,
              ),
              unselectedLabelColor: Colors.grey[400],
              indicatorSize: TabBarIndicatorSize.label,
              indicatorWeight: 3,
              indicatorColor: Colors.blue[700],
            ),
          ),
          const SizedBox(width: 25),
        ],
      ),
      actions: showSearch
          ? [
              IconButton(
                icon: Icon(Icons.search, color: Colors.blue[700]),
                onPressed: () {
                  // 跳转到搜索页面
                },
              )
            ]
          : [],
    );
  }
}
