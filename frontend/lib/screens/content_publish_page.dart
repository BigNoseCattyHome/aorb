import 'dart:ui';
import 'package:aorb/utils/ip_locator.dart';
import 'package:flutter/material.dart';
import 'dart:math';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/generated/poll.pb.dart';
import 'package:shared_preferences/shared_preferences.dart';

class ContentPublishPage extends StatefulWidget {
  const ContentPublishPage({super.key});

  @override
  ContentPublishPageState createState() => ContentPublishPageState();
}

class ContentPublishPageState extends State<ContentPublishPage>
    with SingleTickerProviderStateMixin {
  // 使用Map来存储选项,包含id和文本
  List<Map<String, String>> options = [
    {'id': '1', 'text': '选项1'},
    {'id': '2', 'text': '选项2'}
  ];
  bool isPublic = true;
  TextEditingController titleController = TextEditingController();
  TextEditingController descriptionController = TextEditingController();
  late AnimationController _animationController;

  String ipaddress = '';
  String username = '';

  // 存储背景圆圈的信息
  late List<Map<String, dynamic>> backgroundCircles;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 500),
      vsync: this,
    );

    // 在初始化时生成背景圆圈
    _generateBackgroundCircles();

    // 初始化IP地址和用户名
    _initializeData();
  }

  // 生成背景圆圈的方法
  void _generateBackgroundCircles() {
    final random = Random();
    backgroundCircles = List.generate(3, (index) {
      return {
        'size': random.nextDouble() * 300 + 100,
        'left': random.nextDouble(),
        'top': random.nextDouble(),
        'color': random.nextBool() ? Colors.red : Colors.blue,
      };
    });
  }

  // 新增方法来处理异步初始化
  Future<void> _initializeData() async {
    ipaddress = await IPLocationUtil.getProvince();
    SharedPreferences prefs = await SharedPreferences.getInstance();
    username = prefs.getString('username') ?? '';
    setState(() {}); // 更新状态以反映新的数据
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Stack(
        children: [
          _buildBackground(context),
          SafeArea(
            child: Column(
              children: [
                _buildAppBar(),
                Expanded(
                  child: ListView(
                    padding: const EdgeInsets.all(16),
                    children: [
                      _buildInputBox(),
                      const SizedBox(height: 20),
                      _buildOptionsGrid(),
                      const SizedBox(height: 20),
                      _buildAddOptionButton(),
                    ],
                  ),
                ),
              ],
            ),
          ),
          Positioned(
            right: 16,
            bottom: 16,
            child: _buildVisibilityToggle(),
          ),
        ],
      ),
    );
  }

  Widget _buildBackground(BuildContext context) {
    final screenSize = MediaQuery.of(context).size;
    return Stack(
      children: [
        ...backgroundCircles.map((circle) => Positioned(
              left: circle['left'] * screenSize.width,
              top: circle['top'] * screenSize.height,
              child: Container(
                width: circle['size'],
                height: circle['size'],
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: (circle['color'] as Color).withOpacity(0.7),
                ),
              ),
            )),
        BackdropFilter(
          filter: ImageFilter.blur(sigmaX: 80, sigmaY: 80),
          child: Container(
            color: Colors.transparent,
            child: Container(
              decoration: BoxDecoration(
                gradient: LinearGradient(
                  begin: Alignment.topLeft,
                  end: Alignment.bottomRight,
                  colors: [
                    Colors.blue.withOpacity(0.3),
                    Colors.red.withOpacity(0.3),
                  ],
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }

  Widget _buildAppBar() {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          IconButton(
            icon: const Icon(Icons.close),
            onPressed: () => Navigator.pop(context),
          ),
          TextButton(
            onPressed: () {
              var poll = Poll()
                ..title = titleController.text
                ..content = descriptionController.text
                ..pollType = isPublic ? "public" : "private"
                ..options.addAll(options.map((option) => option['text']!))
                ..ipaddress = ipaddress
                ..username = username;
              var request = CreatePollRequest()..poll = poll;
              PollService().createPoll(request).then((response) {
                if (response.statusCode == 0) {
                  Navigator.pop(context);
                } else {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(response.statusMsg),
                      duration: const Duration(seconds: 2),
                    ),
                  );
                }
              });
            },
            style: TextButton.styleFrom(
              backgroundColor: const Color(0xFF1967DD),
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: const Text(
              '发布',
              style: TextStyle(
                color: Colors.white,
                fontSize: 16,
                fontFamily: 'Inter',
                fontWeight: FontWeight.w700,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildInputBox() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.grey.withOpacity(0.2),
        borderRadius: BorderRadius.circular(15),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          TextField(
            controller: titleController,
            style: const TextStyle(
              color: Colors.black,
              fontSize: 16,
              fontFamily: 'Inter',
              fontWeight: FontWeight.w700,
            ),
            decoration: const InputDecoration(
              hintText: '想要问点大伙什么好呢？',
              border: InputBorder.none,
            ),
          ),
          const Divider(color: Colors.black),
          TextField(
            controller: descriptionController,
            style: const TextStyle(
              color: Colors.black,
              fontSize: 13,
              fontFamily: 'Inter',
              fontWeight: FontWeight.w400,
            ),
            decoration: const InputDecoration(
              hintText: '在这里补充说明一下吧...',
              border: InputBorder.none,
            ),
            maxLines: null,
          ),
        ],
      ),
    );
  }

  Widget _buildOptionsGrid() {
    return GridView.builder(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 2,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: options.length,
      itemBuilder: (context, index) {
        return _buildOptionButton(options[index], index);
      },
    );
  }

  Widget _buildOptionButton(Map<String, String> option, int index) {
    return Dismissible(
      key: Key(option['id']!),
      direction: DismissDirection.endToStart,
      onDismissed: (direction) {
        setState(() {
          options.removeWhere((item) => item['id'] == option['id']);
        });
      },
      background: Container(
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: 20),
        color: Colors.transparent,
        child: const Icon(Icons.delete_forever_rounded,
            color: Colors.white, size: 30),
      ),
      child: GestureDetector(
        onTap: () {
          _editOption(option['id']!);
        },
        child: Container(
          decoration: BoxDecoration(
            color: index % 2 == 0
                ? const Color(0xCCBE2F2F)
                : const Color(0xCC4376DB),
            borderRadius: BorderRadius.circular(20),
          ),
          child: Center(
            child: Text(
              option['text']!,
              style: const TextStyle(
                color: Colors.white,
                fontSize: 24,
                fontFamily: 'Inter',
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _editOption(String id) {
    final optionIndex = options.indexWhere((item) => item['id'] == id);
    if (optionIndex == -1) return;

    showDialog(
      context: context,
      builder: (context) {
        String newOptionText = options[optionIndex]['text']!;
        return AlertDialog(
          title: const Text('编辑选项'),
          content: TextField(
            onChanged: (value) {
              newOptionText = value;
            },
            controller:
                TextEditingController(text: options[optionIndex]['text']),
          ),
          actions: [
            TextButton(
              child: const Text('取消'),
              onPressed: () => Navigator.pop(context),
            ),
            TextButton(
              child: const Text('确定'),
              onPressed: () {
                setState(() {
                  options[optionIndex]['text'] = newOptionText;
                });
                Navigator.pop(context);
              },
            ),
          ],
        );
      },
    );
  }

  Widget _buildAddOptionButton() {
    return GestureDetector(
      onTap: () {
        setState(() {
          options.add({
            'id': DateTime.now().millisecondsSinceEpoch.toString(),
            'text': '新选项'
          });
        });
      },
      child: Container(
        width: 50,
        height: 50,
        decoration: const BoxDecoration(
          shape: BoxShape.circle,
          gradient: LinearGradient(
            colors: [Color(0xFFD83333), Color(0xFF2F88FF)],
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
          ),
        ),
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildVisibilityToggle() {
    return GestureDetector(
      onTap: () {
        setState(() {
          isPublic = !isPublic;
          if (_animationController.status == AnimationStatus.completed) {
            _animationController.reverse();
          } else {
            _animationController.forward();
          }
        });
      },
      child: Container(
        padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
        decoration: BoxDecoration(
          gradient: const LinearGradient(
            colors: [Color(0xFFFF7D7D), Color(0xFF4695FF)],
            begin: Alignment.centerLeft,
            end: Alignment.centerRight,
          ),
          borderRadius: BorderRadius.circular(10),
        ),
        child: Text(
          isPublic ? '公开' : '私密',
          style: const TextStyle(
            color: Colors.white,
            fontSize: 15,
            fontFamily: 'Inter',
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
    );
  }
}
