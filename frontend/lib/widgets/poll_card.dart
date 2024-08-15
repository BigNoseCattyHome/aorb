import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/screens/login_prompt_page.dart';
import 'package:aorb/screens/poll_detail_page.dart';
import 'package:aorb/utils/auth_provider.dart';
import 'package:aorb/utils/color_analyzer.dart';
import 'package:aorb/utils/container_decoration.dart';
import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:aorb/utils/time.dart';
import 'package:provider/provider.dart';

class PollCard extends StatefulWidget {
  final String title; // 投票的标题
  final String content; // 投票的内容说明
  final List<String> options; // 投票的选项
  final List<double> votePercentage; // 用户投票的百分比，根据 options_count 计算
  final int voteCount; // 投票总数
  final Timestamp time; // 投票创建的时间
  final String username; // 用户名，传递给下一个组件
  final String avatar; // 用户头像，根据username查询，在父组件完成的
  final String nickname; // 用户昵称，根据username查询
  final String pollId; // 就是 poll_uuid
  final String userId; // 用户id，根据username查询
  final String backgroundImage; // 背景图片，可以是纯色、网络图片或渐变，根据username查询
  final String selectedOption; // 用户选择的选项,-1代表没有投票

  const PollCard({
    Key? key,
    required this.title,
    required this.content,
    required this.options,
    required this.voteCount,
    required this.time,
    required this.username,
    required this.avatar,
    required this.nickname,
    required this.userId,
    required this.pollId,
    required this.backgroundImage,
    required this.votePercentage,
    required this.selectedOption,
  }) : super(key: key);

  @override
  PollCardState createState() => PollCardState();
}

class PollCardState extends State<PollCard> {
  String _selectedOption = "";
  Color? _textColor;
  var logger = getLogger();

  @override
  void initState() {
    super.initState();
    _selectedOption = widget.selectedOption; // 目前是从父组件传递过来的
    _initializeTextColor();
  }

  Future<void> _initializeTextColor() async {
    _textColor = await ColorAnalyzer.getTextColor(widget.backgroundImage);
    if (mounted) {
      setState(() {});
    }
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
                    Text(formatTimestamp(widget.time, "发布于"),
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
                              color:
                                  (_textColor ?? Colors.black).withOpacity(0.8),
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

  Widget _buildOptionButtons(List<String> optionsText,
      List<double> votePercentage, String selectedOption) {
    // 颜色数组
    // 左边的选项用红色，右边的选项用蓝色
    // 投票后未选中的投片用灰色，选中该的选项保留原来的颜色
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

    bool hasUserSelected = selectedOption.isNotEmpty;

    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(optionsText.length, (i) {
        bool isSelected = selectedOption == optionsText[i];
        final authProvider = context.read<AuthProvider>();
        return Expanded(
          child: Stack(
            children: [
              LayoutBuilder(
                builder: (BuildContext context, BoxConstraints constraints) {
                  double containerWidth = constraints.maxWidth;
                  return ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      backgroundColor: hasUserSelected
                          ? (isSelected
                              ? colorBackground[i + 1]
                              : colorBackground[0])
                          : colorBackground[i + 1],
                      foregroundColor: Colors.white,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(10),
                      ),
                      minimumSize: const Size(130, 50),
                      padding: EdgeInsets.zero,
                    ),
                    child: Stack(
                      alignment: Alignment.center,
                      children: [
                        if (hasUserSelected)
                          Align(
                            alignment: Alignment.centerLeft,
                            child: Container(
                              height: 50,
                              width: votePercentage[i] * containerWidth,
                              decoration: BoxDecoration(
                                color: isSelected
                                    ? colorPercents[i + 1]
                                    : colorPercents[0],
                                borderRadius: BorderRadius.circular(10),
                              ),
                            ),
                          ),
                        Align(
                          alignment: Alignment.center,
                          child: Text(optionsText[i]),
                        ),
                      ],
                    ),
                    onPressed: () {
                      authProvider.isLoggedIn
                          ? Navigator.push(
                              context,
                              MaterialPageRoute(
                                builder: (context) => PollDetailPage(
                                  userId: widget.userId,
                                  pollId: widget.pollId,
                                  username: widget.username,
                                ),
                              ),
                            )
                          : const LoginPromptPage();
                    },
                  );
                },
              )
            ],
          ),
        );
      }),
    );
  }

  // 点击跳转到内容详情页面
  void onTapContent() {
    final authProvider = context.read<AuthProvider>();
    authProvider.isLoggedIn
        ? Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => PollDetailPage(
                userId: widget.userId,
                pollId: widget.pollId,
                username: widget.username,
              ),
            ),
          )
        : const LoginPromptPage();
  }
}
