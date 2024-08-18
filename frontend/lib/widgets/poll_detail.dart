import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/services/vote_service.dart';
import 'package:aorb/utils/color_analyzer.dart';
import 'package:aorb/utils/container_decoration.dart';
import 'package:aorb/utils/time.dart';
import 'package:flutter/material.dart';

class PollDetail extends StatefulWidget {
  final String title;
  final String content;
  final List<String> options;
  final List<double> votePercentage;
  final Timestamp time;
  final String ipaddress;
  final String selectedOption;
  final String pollId;
  final String username;
  final String bgpic;
  final String currentUser;

  const PollDetail({
    Key? key,
    required this.title,
    required this.content,
    required this.options,
    required this.votePercentage,
    required this.time,
    required this.ipaddress,
    required this.selectedOption,
    required this.pollId,
    required this.username,
    required this.bgpic,
    required this.currentUser,
  }) : super(key: key);

  @override
  PollDetailState createState() => PollDetailState();
}

class PollDetailState extends State<PollDetail> {
  String _selectedOption = "";
  Color _textColor = Colors.black; // 默认文字颜色

  @override
  void initState() {
    super.initState();
    fetchColor();
    _selectedOption = widget.selectedOption;
  }

  void fetchColor() async {
    _textColor = await ColorAnalyzer.getTextColor(widget.bgpic);
    if (mounted) {
      setState(() {});
    }
  }

  void _selectOption(int index) {
    setState(() {
      _selectedOption = widget.options[index];
    });
    VoteService().createVote(
      widget.pollId,
      widget.currentUser,
      widget.options[index],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 15.0),
      decoration: createBackgroundDecoration(widget.bgpic),
      child: Padding(
        padding: const EdgeInsets.all(25.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              widget.title,
              style: TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: _textColor,
              ),
            ),
            const SizedBox(height: 16),
            Text(
              widget.content,
              style: TextStyle(
                fontSize: 16,
                color: _textColor,
              ),
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _buildOptionButton(0, const Color(0xFFBE3030)),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: _buildOptionButton(1, const Color(0xFF4376DB)),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Text(
                  formatTimestamp(widget.time, "发布于"),
                  style: TextStyle(
                    fontSize: 14,
                    color: _textColor.withOpacity(0.6),
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  widget.ipaddress,
                  style: TextStyle(
                    fontSize: 14,
                    color: _textColor.withOpacity(0.6),
                  ),
                )
              ],
            )
          ],
        ),
      ),
    );
  }

  Widget _buildOptionButton(int index, Color color) {
    return GestureDetector(
      onTap: () => _selectOption(index),
      child: Container(
        padding: const EdgeInsets.symmetric(vertical: 12),
        decoration: BoxDecoration(
          color: _selectedOption == widget.options[index]
              ? color
              : color.withOpacity(0.6),
          borderRadius: BorderRadius.circular(8),
        ),
        child: Center(
          child: Text(
            widget.options[index],
            style: const TextStyle(
              color: Colors.white,
              fontSize: 16,
            ),
          ),
        ),
      ),
    );
  }
}
