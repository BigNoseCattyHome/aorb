import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:intl/intl.dart';

class PollCard extends StatefulWidget {
  final String title; // 投票的标题
  final String content; // 投票的内容说明
  final List<String> options; // 投票的选项
  final List<double> votePercentage; // 用户投票的百分比，根据 options_count 计算
  final int voteCount; // 投票总数
  final Timestamp time; // 投票创建的时间
  final String avatar; // 用户头像，根据username查询
  final String nickname; // 用户昵称，根据username查询
  final String pollId; // 就是 poll_uuid
  final String userId; // 用户id，根据username查询
  final String backgroundImage; // 背景图片，可以是纯色、网络图片或渐变，根据username查询
  final int selectedOption; // 用户选择的选项,-1代表没有投票

  const PollCard({
    Key? key,
    required this.title,
    required this.content,
    required this.options,
    required this.voteCount,
    required this.time,
    required this.avatar,
    required this.nickname,
    required this.userId,
    required this.pollId,
    required this.backgroundImage,
    required this.votePercentage,
    this.selectedOption = -1,
  }) : super(key: key);

  @override
  PollCardState createState() => PollCardState();
}

class PollCardState extends State<PollCard> {
  int _selectedOption = -1;
  late Color _textColor;

  @override
  void initState() {
    super.initState();
    // widget就是当前State对象关联的Widget对象
    _selectedOption = widget.selectedOption;
    _textColor = _getTextColor(widget.backgroundImage);
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16.0), // 增加外部间距
      child: GestureDetector(
          onTap: onTapContent, // 将onTapContent方法绑定到点击事件上
          child: Container(
            // decoration是容器的装饰，这里设置了背景图片和圆角
            decoration: createBackgroundDecoration(widget.backgroundImage),

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
                      style: TextStyle(
                          color: _textColor, fontWeight: FontWeight.bold),
                    ),
                    const Spacer(), // spacer占位部件
                    // time
                    Text(_formatTimestamp(widget.time),
                        style: TextStyle(
                          color: _textColor,
                        )),
                  ],
                ),
                const SizedBox(height: 16),
                Stack(
                  children: [
                    Positioned(
                      top: 3,
                      left: 0,
                      child: SvgPicture.asset('images/comments.svg'),
                    ),
                    const SizedBox(width: 8),
                    Padding(
                      padding:
                          const EdgeInsets.only(left: 60), // 根据你的SVG图标的大小调整这个值
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          // title
                          Text(
                            widget.title,
                            style: TextStyle(
                              color: _textColor,
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 8),
                          // content
                          Text(
                            widget.content,
                            style: TextStyle(
                              color: _textColor.withOpacity(0.8),
                              fontSize: 10,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),

                const SizedBox(height: 16),
                // options
                _buildOptionButtons(
                    widget.options, widget.votePercentage, _selectedOption),
              ],
            ),
          )),
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
        // Expanded是一个占满剩余空间的部件
        return Expanded(
            child: Stack(children: [
          // 获取ElevatedButton部件的宽度
          LayoutBuilder(
              builder: (BuildContext context, BoxConstraints constraints) {
            // 获取Container的最大宽度
            double containerWidth = constraints.maxWidth;
            // ^ 获取Container的最大高度
            // double containerHeight = constraints.maxHeight;
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
                minimumSize: const Size(130, 50),
                padding: EdgeInsets.zero, // 设置内边距为0
              ),

              child: Stack(
                alignment: Alignment.center, // 设置Stack的对齐方式为居中
                children: [
                  if (_selectedOption != -1)
                    Align(
                      alignment: Alignment.centerLeft,
                      child: Container(
                        height: 40, // ! 好像containerHeight的数值有问题，所以先用50测试
                        // 宽度等于父组件的长度*这个option所占的百分比
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

              // 单次点击投票
              onPressed: () {
                setState(() {
                  _selectedOption = i;
                  // TODO 把选择的结果发送到服务端
                });
              },

              // 长按取消
              onLongPress: () {
                setState(() {
                  _selectedOption = -1;
                  // TODO 把选择的结果发送到服务端
                });
              },
            );
          })
        ]));
      }),
    );
  }

  // 方法：根据backgroundImage创建背景装饰
  BoxDecoration createBackgroundDecoration(String backgroundImage) {
    if (backgroundImage.startsWith('0x')) {
      // 纯色背景
      int colorValue = int.parse(backgroundImage.substring(2), radix: 16);
      return BoxDecoration(
        color: Color(colorValue),
        borderRadius: BorderRadius.circular(10),
      );
    } else if (backgroundImage.startsWith('http://') ||
        backgroundImage.startsWith('https://')) {
      // 网络图片背景
      return BoxDecoration(
        image: DecorationImage(
          image: NetworkImage(backgroundImage),
          fit: BoxFit.cover,
        ),
        borderRadius: BorderRadius.circular(10),
      );
    } else if (backgroundImage.startsWith('gradient:')) {
      // 渐变背景
      List<String> colorStrings =
          backgroundImage.substring('gradient:'.length).split(',');
      List<Color> colors = colorStrings
          .map((colorString) =>
              Color(int.parse(colorString.substring(2), radix: 16)))
          .toList();
      return BoxDecoration(
        gradient: LinearGradient(
          colors: colors,
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(10),
      );
    } else {
      // 默认纯色背景
      return BoxDecoration(
        color: Colors.blue[900],
        borderRadius: BorderRadius.circular(10),
      );
    }
  }

  // 点击跳转到内容详情页面
  void onTapContent() {
    Navigator.pushNamed(
      context,
      '/poll_content',
      arguments: {'postId': widget.pollId, 'userId': widget.userId},
    );
  }

  // 新增方法：格式化时间戳
  String _formatTimestamp(Timestamp timestamp) {
    final now = DateTime.now();
    final dateTime = timestamp.toDateTime();
    final difference = now.difference(dateTime);

    if (difference.inMinutes < 60) {
      return '发布于${difference.inMinutes}分钟前';
    } else if (difference.inHours < 24) {
      return '发布于${difference.inHours}小时前';
    } else if (difference.inDays < 2) {
      return '发布于昨天${DateFormat('HH:mm').format(dateTime)}';
    } else if (difference.inDays < 7) {
      return '发布于${difference.inDays}天前';
    } else if (difference.inDays < 30) {
      return '发布于${(difference.inDays / 7).floor()}周前';
    } else if (difference.inDays < 365) {
      return '发布于${(difference.inDays / 30).floor()}个月前';
    } else {
      return '发布于${DateFormat('yyyy年MM月dd日').format(dateTime)}';
    }
  }

// 新增方法：根据背景决定文字颜色
  Color _getTextColor(String backgroundImage) {
    if (backgroundImage.startsWith('0x')) {
      // 纯色背景
      int colorValue = int.parse(backgroundImage.substring(2), radix: 16);
      Color backgroundColor = Color(colorValue);
      return _isLightColor(backgroundColor) ? Colors.black : Colors.white;
    } else if (backgroundImage.startsWith('gradient:')) {
      // 渐变背景，我们取第一个颜色作为参考
      List<String> colorStrings =
          backgroundImage.substring('gradient:'.length).split(',');
      Color firstColor =
          Color(int.parse(colorStrings[0].substring(2), radix: 16));
      return _isLightColor(firstColor) ? Colors.black : Colors.white;
    } else {
      // 对于图片背景，我们默认使用白色文字
      return Colors.white;
    }
  }

  // 新增方法：判断颜色是否为浅色
  bool _isLightColor(Color color) {
    // 使用相对亮度公式
    return (0.299 * color.red + 0.587 * color.green + 0.114 * color.blue) /
            255 >
        0.5;
  }
}