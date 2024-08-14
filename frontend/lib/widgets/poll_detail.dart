import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/services/vote_service.dart';
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

  @override
  void initState() {
    super.initState();
    _selectedOption = widget.selectedOption;
  }

  void _selectOption(int index) {
    setState(() {
      _selectedOption = widget.options[index];
    });
    VoteService().createVote(
      widget.pollId,
      widget.username,
      widget.options[index],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        image: DecorationImage(
          image: NetworkImage(widget.bgpic),
          fit: BoxFit.cover,
        ),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              widget.title,
              style: const TextStyle(
                fontSize: 24,
                fontWeight: FontWeight.bold,
                color: Colors.black,
              ),
            ),
            const SizedBox(height: 16),
            Text(
              widget.content,
              style: const TextStyle(
                fontSize: 16,
                color: Colors.black,
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
                  style: const TextStyle(
                    fontSize: 14,
                    color: Colors.grey,
                  ),
                ),
                const SizedBox(width: 8),
                Text(
                  widget.ipaddress,
                  style: const TextStyle(
                    fontSize: 14,
                    color: Colors.grey,
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
          color: _selectedOption == widget.options[index] ? color : color.withOpacity(0.6),
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