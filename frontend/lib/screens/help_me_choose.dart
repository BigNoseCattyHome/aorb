import 'package:flutter/material.dart';
import 'system_choice.dart';
import 'public_vote.dart';
import 'private_vote.dart';
import 'anonymous_vote.dart';

class HelpMeChoose extends StatelessWidget {
  void _showModalBottomSheet(BuildContext context) {
    showModalBottomSheet(
      context: context,
      builder: (BuildContext context) {
        return Container(
          height: 250,
          child: ListView(
            children: <Widget>[
              ListTile(
                leading: Icon(Icons.computer),
                title: Text('系统选择'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => SystemChoice()));
                },
              ),
              ListTile(
                leading: Icon(Icons.public),
                title: Text('公开投票'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => PublicVote()));
                },
              ),
              ListTile(
                leading: Icon(Icons.lock_outline),
                title: Text('私密投票'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => PrivateVote()));
                },
              ),
              ListTile(
                leading: Icon(Icons.visibility_off),
                title: Text('匿名投票'),
                onTap: () {
                  Navigator.pop(context);
                  Navigator.push(context, MaterialPageRoute(builder: (_) => AnonymousVote()));
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
        title: Text('帮我选'),
      ),
      body: Center(
        child: ElevatedButton(
          onPressed: () => _showModalBottomSheet(context),
          child: Text('显示选项'),
        ),
      ),
    );
  }
}
