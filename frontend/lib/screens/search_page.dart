import 'package:flutter/material.dart';

class SearchPage extends StatefulWidget {
  final TabController tabController;

  const SearchPage({Key? key, required this.tabController}) : super(key: key);

  @override
  SearchPageState createState() => SearchPageState();
}

class SearchPageState extends State<SearchPage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  late TextEditingController _searchController;

  @override
  void initState() {
    super.initState();
    _tabController = widget.tabController;
    _searchController = TextEditingController();
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: TextField(
          controller: _searchController,
          decoration: const InputDecoration(
            hintText: 'Search',
            border: InputBorder.none,
          ),
          onSubmitted: (value) {
            // 搜索
          },
        ),
        actions: [
          IconButton(
            onPressed: () {
              // 清空搜索框
              _searchController.clear();
            },
            icon: const Icon(Icons.clear),
          ),
        ],
      ),
      body: TabBarView(
        controller: _tabController,
        children: const [
          // 搜索结果
        ],
      ),
    );
  }
}
