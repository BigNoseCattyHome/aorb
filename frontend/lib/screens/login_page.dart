import 'package:flutter/material.dart';
import 'package:aorb/screens/register_page.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});
  
  @override
  LoginPageState createState() => LoginPageState();
}

class LoginPageState extends State<LoginPage> {
  bool _obscureText = true;
  void _toggleObscureText() {
    setState(() {
      _obscureText = !_obscureText;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Padding(
        padding: const EdgeInsets.all(32.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: <Widget>[
            const Text(
              '欢迎来到Aorb',
              style: TextStyle(
                fontSize: 40,
                fontWeight: FontWeight.bold,
              ),
              textAlign: TextAlign.left, // 左对齐
            ),
            const SizedBox(height: 8),
            const Text(
              '登录账户解锁更多功能~',
              style: TextStyle(
                fontSize: 20,
                color: Colors.black87,
                fontWeight: FontWeight.bold,
              ),
              textAlign: TextAlign.left, // 左对齐
            ),
            const SizedBox(height: 32),
            CircleAvatar(
              radius: 50,
              backgroundColor: Colors.grey[300],
              child: Icon(
                Icons.person,
                size: 50,
                color: Colors.grey[600],
              ),
            ),

            const SizedBox(height: 32),

            TextFormField(
              decoration: const InputDecoration(
                labelText: '用户名',
                border: InputBorder.none,
                enabledBorder: UnderlineInputBorder(
                  borderSide: BorderSide(color: Colors.grey),
                ),
                focusedBorder: UnderlineInputBorder(
                  borderSide: BorderSide(color: Colors.blue),
                ),
                prefixIcon: Icon(Icons.person),
              ),
            ),

            const SizedBox(height: 16),
            TextFormField(
              obscureText: _obscureText,
              decoration: InputDecoration(
                labelText: '密码',
                border: InputBorder.none,
                enabledBorder: UnderlineInputBorder(
                  borderSide: BorderSide(color: Colors.grey),
                ),
                focusedBorder: UnderlineInputBorder(
                  borderSide: BorderSide(color: Colors.blue),
                ),
                prefixIcon: Icon(Icons.key),
                suffixIcon: IconButton(
                  icon: Icon(
                    _obscureText ? Icons.visibility : Icons.visibility_off,
                  ),
                  onPressed: _toggleObscureText,
                ),
              ),
            ),
            const SizedBox(height: 24),
            // 圆形的按钮，内含右边的箭头，表示登录的按钮
            MaterialButton(
              color: Colors.blue,
              textColor: Colors.white,
              padding: const EdgeInsets.all(16),
              shape: const CircleBorder(),
              child: const Icon(Icons.arrow_forward),
              onPressed: () {
                // TODO: 登录逻辑
              },
            ),
            const SizedBox(height: 16),
            TextButton(
              onPressed: () {
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (context) => RegisterPage(),
                  ),
                );
              },
              child: const Text('没有账户？注册'),
            ),
          ],
        ),
      ),
    );
  }
}
