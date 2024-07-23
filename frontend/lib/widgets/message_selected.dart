import 'package:flutter/material.dart';

class MessageSelected extends StatelessWidget {
  final String message;
  final String time;
  final List<String> options;

  const MessageSelected({
    Key? key,
    required this.message,
    required this.time,
    required this.options,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(vertical: 8.0),
      padding: const EdgeInsets.all(16.0),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12.0),
        boxShadow: const [
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
          Text(message, style: const TextStyle(color: Colors.black)),
          const SizedBox(height: 8.0),
          Wrap(
            spacing: 8.0,
            children: options.map((option) => Chip(label: Text(option))).toList(),
          ),
          const SizedBox(height: 8.0),
          Text(time, style: const TextStyle(color: Colors.grey)),
        ],
      ),
    );
  }
}
