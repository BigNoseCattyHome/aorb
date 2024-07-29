import 'package:aorb/conf/config.dart';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SettingsPage extends StatelessWidget {
  SettingsPage({super.key});
  final logger = getLogger();
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Center(
        child: ElevatedButton(
          onPressed: () async {
            final prefs = await SharedPreferences.getInstance();
            await prefs.clear();
            logger.i('All data cleared.');
          },
          child: Text('Clear Preferences'),
        ),
      ),
    );
  }
}
