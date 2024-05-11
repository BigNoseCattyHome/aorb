import 'package:flutter/material.dart';
import 'system_choice.dart';
import 'public_vote.dart';
import 'private_vote.dart';
import 'anonymous_vote.dart';

class HelpMeChoose extends StatelessWidget {
  const HelpMeChoose({super.key});

  void _showModalBottomSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      builder: (BuildContext context) {
        return SizedBox(
          height: 250,
          child: ListView(
            children: <Widget>[
              ListTile(
                leading: const Icon(Icons.computer),
                title: const Text('系统选择'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => const SystemChoice()));
                },
              ),
              ListTile(
                leading: const Icon(Icons.public),
                title: const Text('公开投票'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => const PublicVote()));
                },
              ),
              ListTile(
                leading: const Icon(Icons.lock_outline),
                title: const Text('私密投票'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => const PrivateVote()));
                },
              ),
              ListTile(
                leading: const Icon(Icons.visibility_off),
                title: const Text('匿名投票'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => const AnonymousVote()));
                },
              ),
            ],
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('帮我选'),
      ),
      body: Center(
        child: ElevatedButton(
          onPressed: () => _showModalBottomSheet(context),
          child: const Text('显示选项'),
        ),
      ),
    );
  }
}
