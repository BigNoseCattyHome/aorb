import 'package:flutter/material.dart';
import 'package:intl/intl.dart';

class PollDetail extends StatefulWidget {
  final String title;
  final String content;
  final List<String> options;
  final List<double> votePercentage;
  final DateTime time;
  final String ipaddress;

  const PollDetail({
    Key? key,
    required this.title,
    required this.content,
    required this.options,
    required this.votePercentage,
    required this.time,
    required this.ipaddress,
  }) : super(key: key);

  @override
  PollDetailState createState() => PollDetailState();
}

class PollDetailState extends State<PollDetail> {
  int _selectedOption = -1;

  @override
  void initState() {
    super.initState();
    _selectedOption = -1;
  }

  void _selectOption(int index) {
    setState(() {
      _selectedOption = index;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(16.0),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            widget.title,
            style: TextStyle(
              fontSize: 24,
              fontWeight: FontWeight.bold,
            ),
          ),
          SizedBox(height: 16),
          Text(
            widget.content,
            style: TextStyle(
              fontSize: 16,
            ),
          ),
          SizedBox(height: 16),
          for (int i = 0; i < widget.options.length; i++)
            GestureDetector(
              onTap: () => _selectOption(i),
              child: Container(
                margin: EdgeInsets.only(bottom: 8),
                padding: EdgeInsets.symmetric(vertical: 12),
                decoration: BoxDecoration(
                  color: _selectedOption == i ? Colors.red : Colors.grey,
                  borderRadius: BorderRadius.circular(8),
                ),
                child: Center(
                  child: Text(
                    widget.options[i],
                    style: TextStyle(
                      color: Colors.white,
                      fontSize: 16,
                    ),
                  ),
                ),
              ),
            ),
          SizedBox(height: 16),
          Text(
            '发布于 ${DateFormat('H小时').format(widget.time)}前  ${widget.ipaddress}',
            style: TextStyle(
              fontSize: 14,
              color: Colors.grey,
            ),
          ),
        ],
      ),
    );
  }
}
