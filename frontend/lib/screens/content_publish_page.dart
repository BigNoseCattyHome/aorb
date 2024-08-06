import 'dart:ui';
import 'package:flutter/material.dart';
import 'dart:math';

class ContentPublishPage extends StatefulWidget {
  const ContentPublishPage({super.key});

  @override
  ContentPublishPageState createState() => ContentPublishPageState();
}

class ContentPublishPageState extends State<ContentPublishPage>
    with SingleTickerProviderStateMixin {
  List<String> options = ['选项1', '选项2'];
  bool isPublic = true;
  TextEditingController titleController = TextEditingController();
  TextEditingController descriptionController = TextEditingController();
  late AnimationController _animationController;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 500),
      vsync: this,
    );
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
    return Stack(
      children: [
        ..._buildBackgroundCircles(context),
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

  List<Widget> _buildBackgroundCircles(BuildContext context) {
    final screenSize = MediaQuery.of(context).size;
    final random = Random();
    return List.generate(3, (index) {
      final size = random.nextDouble() * 300 + 100;
      final left = random.nextDouble() * screenSize.width;
      final top = random.nextDouble() * screenSize.height;
      final color = random.nextBool()
          ? Colors.red.withOpacity(0.7)
          : Colors.blue.withOpacity(0.7);

      return Positioned(
        left: left,
        top: top,
        child: Container(
          width: size,
          height: size,
          decoration: BoxDecoration(
            shape: BoxShape.circle,
            color: color,
          ),
        ),
      );
    });
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
              // 实现发布功能
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

  Widget _buildOptionButton(String option, int index) {
    return Dismissible(
      key: Key(option),
      direction: DismissDirection.endToStart,
      onDismissed: (direction) {
        setState(() {
          options.removeAt(index);
        });
      },
      background: Container(
        alignment: Alignment.centerRight,
        padding: const EdgeInsets.only(right: 20),
        color: Colors.transparent,
        child:
            const Icon(Icons.delete_forever_rounded, color: Colors.white, size: 30),
      ),
      child: GestureDetector(
        onTap: () {
          _editOption(index);
        },
        child: Container(
          decoration: BoxDecoration(
            color: index % 2 == 0 ? const Color(0xCCBE2F2F) : const Color(0xCC4376DB),
            borderRadius: BorderRadius.circular(20),
          ),
          child: Center(
            child: Text(
              option,
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

  void _editOption(int index) {
    showDialog(
      context: context,
      builder: (context) {
        String newOption = options[index];
        return AlertDialog(
          title: const Text('编辑选项'),
          content: TextField(
            onChanged: (value) {
              newOption = value;
            },
            controller: TextEditingController(text: options[index]),
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
                  options[index] = newOption;
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
          options.add('新选项');
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
