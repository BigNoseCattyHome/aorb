import 'dart:math';

import 'package:aorb/generated/auth.pbgrpc.dart';
import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:flutter/material.dart';
import 'package:aorb/screens/register_page.dart';
import 'package:aorb/conf/config.dart';
import 'package:aorb/services/auth_service.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  LoginPageState createState() => LoginPageState();
}

class LoginPageState extends State<LoginPage> {
  final _usernameController = TextEditingController(); //控制输入框
  final _passwordController = TextEditingController();
  final AuthService _authService = AuthService();
  final logger = getLogger();

  bool _obscureText = true; //控制密码可见
  void _toggleObscureText() {
    setState(() {
      _obscureText = !_obscureText;
    });
  }

  void _login() async {
    try {
      logger.d('the username input is: ${_usernameController.text}');
      logger.d('the password input is: ${_passwordController.text}');

      LoginRequest request = LoginRequest()
        ..username = _usernameController.text
        ..password = _passwordController.text
        ..deviceId = 'web'
        ..timestamp = Timestamp.create()
        ..nonce = Random().nextInt(1000000).toString();

      //等待
      final loginResponse = await _authService.login(request);

      logger.i('loginResponse: $loginResponse');

      if (loginResponse.statusCode == 0) {
        // 登录成功后的处理逻辑
        Navigator.pushReplacementNamed(context, '/me');
      } else {
        // 登录失败后的处理逻辑，弹出错误提示
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('登录失败: ${loginResponse.statusMsg}')),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('登录失败: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: ListView(
        padding: const EdgeInsets.all(32.0),
        children: <Widget>[
          Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: <Widget>[
              const SizedBox(height: 120),
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
              const SizedBox(height: 80),
              CircleAvatar(
                radius: 50,
                backgroundColor: Colors.grey[300],
                child: Icon(
                  Icons.person,
                  size: 50,
                  color: Colors.grey[600],
                ),
              ),

              const SizedBox(height: 20),

              TextFormField(
                controller: _usernameController,
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
                controller: _passwordController,
                obscureText: _obscureText,
                decoration: InputDecoration(
                  labelText: '密码',
                  border: InputBorder.none,
                  enabledBorder: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Colors.grey),
                  ),
                  focusedBorder: const UnderlineInputBorder(
                    borderSide: BorderSide(color: Colors.blue),
                  ),
                  prefixIcon: const Icon(Icons.key),
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
                onPressed: _login,
                child: const Icon(Icons.arrow_forward),
              ),
              const SizedBox(height: 16),
              TextButton(
                onPressed: () {
                  Navigator.of(context).push(
                    MaterialPageRoute(
                      builder: (context) => const RegisterPage(),
                    ),
                  );
                },
                child: const Text('没有账户？注册'),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
