import 'package:flutter/material.dart';

class msg_selected extends StatelessWidget {
  final String message;
  final String time;
  final List<String> options;

  const msg_selected({
    Key? key,
    required this.message,
    required this.time,
    required this.options,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: EdgeInsets.symmetric(vertical: 8.0),
      padding: EdgeInsets.all(16.0),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.0),
        boxShadow: [
          BoxShadow(
            color: Colors.black26,
            blurRadius: 6.0,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(message, style: TextStyle(color: Colors.black)),
          SizedBox(height: 8.0),
          Wrap(
            spacing: 8.0,
            children: options.map((option) => Chip(label: Text(option))).toList(),
          ),
          SizedBox(height: 8.0),
          Text(time, style: TextStyle(color: Colors.grey)),
        ],
      ),
    );
  }
}
