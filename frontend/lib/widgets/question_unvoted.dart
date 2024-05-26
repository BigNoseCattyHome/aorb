import 'package:flutter/material.dart';

class QuestionUnvoted extends StatefulWidget {
  final String title;
  final String content;
  final List<String> options;
  final List<double> votePercentage; // 用户投票的百分比
  final int voteCount;
  final String time;
  final String avatar;
  final String nickname;
  final String questionId;
  final String backgroundImage;
  final int selectedOption; // 用户选择的选项,-1代表没有投票

  const QuestionUnvoted({
    Key? key,
    required this.title,
    required this.content,
    required this.options,
    required this.voteCount,
    required this.time,
    required this.avatar,
    required this.nickname,
    required this.questionId,
    required this.backgroundImage,
    required this.votePercentage,
    this.selectedOption = -1,
  }) : super(key: key);

  @override
  QuestionUnvotedState createState() => QuestionUnvotedState();
}

class QuestionUnvotedState extends State<QuestionUnvoted> {
  int _selectedOption = -1;

  @override
  void initState() {
    super.initState();
    // widget就是当前State对象关联的Widget对象
    _selectedOption = widget.selectedOption;
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      // decoration是容器的装饰，这里设置了背景图片和圆角
      decoration: BoxDecoration(
        image: DecorationImage(
          image: NetworkImage(widget.backgroundImage), // network image从网络加载图片
          fit: BoxFit.cover,
        ),
        borderRadius: BorderRadius.circular(10),
      ),

      // padding是容器的内边距，这里设置了上下左右各16像素
      padding: const EdgeInsets.all(16.0),

      // child是容器的子部件，这里是一个Column，包含了问题的标题、内容、选项和投票按钮
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start, // 子部件左对齐
        children: [
          Row(
            children: [
              // avatar
              CircleAvatar(
                backgroundImage: NetworkImage(widget.avatar),
              ),
              const SizedBox(width: 8),
              // nickname
              Text(
                widget.nickname,
                style: const TextStyle(
                    color: Colors.blue, fontWeight: FontWeight.bold),
              ),
              const Spacer(), // spacer占位部件
              // time
              Text(
                widget.time,
                style:
                    const TextStyle(color: Color.fromARGB(255, 255, 255, 255)),
              ),
            ],
          ),
          const SizedBox(height: 16),
          // title
          Text(
            widget.title,
            style: const TextStyle(
                color: Color.fromARGB(255, 255, 255, 255), fontSize: 24),
          ),
          const SizedBox(height: 8),
          // content
          Text(
            widget.content,
            style: const TextStyle(
                color: Color.fromARGB(221, 255, 255, 255), fontSize: 10),
          ),
          const SizedBox(height: 16),
          // options
          _buildOptionButtons(
              widget.options, widget.votePercentage, _selectedOption),
        ],
      ),
    );
  }

  Widget _buildOptionButtons(
      List<String> text, List<double> votePercentage, int selectedOption) {
    // 颜色数组
    List<Color> colorBackground = [
      const Color(0xFF9D9D9D),
      const Color(0xFFBE3030),
      const Color(0xFF4376DB),
    ];
    List<Color> colorPercents = [
      const Color(0xFFB8B8B8),
      const Color(0xFFE55858),
      const Color(0xFF6DB6EB),
    ];

    return Row(
      mainAxisAlignment: MainAxisAlignment.center, // 设置主轴对齐方式为居中
      children: List.generate(text.length, (i) {
        return Expanded(
            // Expanded是一个占满剩余空间的部件
            child:
                // 获取ElevatedButton部件的宽度
                LayoutBuilder(builder:
                    (BuildContext context, BoxConstraints constraints) {
          // 获取Container的最大宽度
          double containerWidth = constraints.maxWidth;
          return ElevatedButton(
            // style是按钮的样式，这里设置了背景颜色、圆角等
            style: ElevatedButton.styleFrom(
              // 设置按钮颜色
              backgroundColor: selectedOption == -1
                  ? colorBackground[i + 1]
                  : (selectedOption == i
                      ? colorBackground[i + 1]
                      : colorBackground[0]),
              // 设置按钮上文本颜色
              foregroundColor: Colors.white,
              // 设置按钮圆角
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(10),
              ),
              minimumSize: const Size(130, 50), // 设置最小高度为 10
            ),
            // 要求当selectedOption更改时，外部调用组件要能够获取到更新
            // 同时组件自己也要更新
            onPressed: () {
              setState(() {
                _selectedOption = i;
              });
            },
            child: Stack(
              children: [
                if (_selectedOption != -1)
                  Positioned(
                    left: 0, // 设置为0以确保Container相对于ElevatedButton左对齐
                    top: 0, // 设置为0以确保Container相对于ElevatedButton顶部对齐
                    child: Container(
                      height: 20, // 设置矩形高度
                      // 宽度等于父组件的长度*这个option所占的百分比
                      // width: votePercentage[i] * 130,
                      width: votePercentage[i] * containerWidth,
                      decoration: BoxDecoration(
                        color: _selectedOption == i
                            ? colorPercents[i + 1] // 使用正确的颜色
                            : colorPercents[0], // 使用默认颜色
                        borderRadius: BorderRadius.circular(10), // 设置圆角半径
                      ),
                    ),
                  ),
                // 使用Align小部件来居中对齐Container
                Align(
                  alignment: Alignment.center,
                  child: Text(text[i]),
                ),
              ],
            ),
          );
        }));
      }),
    );
  }
}
