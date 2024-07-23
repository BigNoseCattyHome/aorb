import 'package:flutter/material.dart';
import 'package:aorb/screens/register_page.dart';
import 'package:aorb/conf/config.dart';
import 'package:aorb/utils/ip_locator.dart';
import 'package:aorb/services/auth_service.dart';

class LoginPage extends StatefulWidget {
  const LoginPage({super.key});
  
  @override
  LoginPageState createState() => LoginPageState();
}

class LoginPageState extends State<LoginPage> {
  
  final _usernameController = TextEditingController(); //控制输入框
  final _passwordController = TextEditingController();
  final AuthService _authService = AuthService(backendHost, backendPort);
  String _province = 'Loading...'; // 用户IP的归属地
  final logger = getLogger();

  bool _obscureText = true; //控制密码可见
  void _toggleObscureText() {
    setState(() {
      _obscureText = !_obscureText;
    });
  }

void _login() async {
  try {
    //等待
    final loginResponse = await _authService.login(
      _usernameController.text,
      _passwordController.text,
      _province  //ip地址，对应login参数，这里只是一个字符串，等待实现获取地址的方法(对应第19行)
    );

    logger.i('loginResponse: $loginResponse');
    
    if (loginResponse.statusCode) {
      // 登录成功后的处理逻辑
      Navigator.pushReplacementNamed(context, '/home');
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
              onPressed: _login,
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
