import 'package:flutter/material.dart';

class DynamicTopBar extends StatefulWidget implements PreferredSizeWidget {
  final List<String> tabs;
  final bool showSearch; // 是否显示搜索按钮

  const DynamicTopBar({
    Key? key, 
    required this.tabs, 
    this.showSearch = true, // 设置默认值为true，变为可选参数
  }) : super(key: key);
  @override
  DynamicTopBarState createState() => DynamicTopBarState();

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);
}

class DynamicTopBarState extends State<DynamicTopBar> with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: widget.tabs.length, vsync: this);
    _tabController.addListener(_handleTabSelection);
  }

  void _handleTabSelection() {
    if (_tabController.indexIsChanging) {
      setState(() {});
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return AppBar(
      elevation: 0,
      backgroundColor: Colors.white,
      leading: IconButton(
        icon: const Icon(Icons.menu, color: Colors.blue),
        onPressed: () {
          // Scaffold.of(context).openDrawer();
        },
      ),
      title: Row(
        children: [
          const SizedBox(width: 25),
          Expanded(
            child: TabBar(
              controller: _tabController,
              tabs: widget.tabs.map((tab) => Tab(text: tab)).toList(),
              labelColor: Colors.blue,
              unselectedLabelColor: Colors.grey,
              indicatorSize: TabBarIndicatorSize.label,
              indicatorWeight: 3,
            ),
          ),
          const SizedBox(width: 25),
        ],
      ),
      actions: widget.showSearch
          ? [
              IconButton(
                icon: const Icon(Icons.search, color: Colors.blue),
                onPressed: () {
                  // 跳转到搜索页面
                },
              )
            ]
          : [],
    );
  }

  Widget buildTabBarView() {
    return TabBarView(
      controller: _tabController,
      children: widget.tabs.map((tab) => Center(child: Text('$tab 内容'))).toList(),
    );
  }
}