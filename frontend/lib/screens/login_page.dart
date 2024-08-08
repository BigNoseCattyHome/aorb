import 'dart:math';
import 'package:flutter/material.dart';
import 'package:aorb/generated/auth.pbgrpc.dart';
import 'package:aorb/generated/google/protobuf/timestamp.pb.dart';
import 'package:aorb/screens/register_page.dart';
import 'package:aorb/conf/config.dart';
import 'package:aorb/services/auth_service.dart';
import 'package:aorb/utils/ip_locator.dart';
import 'package:aorb/utils/constant/err.dart';
import 'package:fluttertoast/fluttertoast.dart';

// 登录页面
class LoginPage extends StatefulWidget {
  const LoginPage({super.key});

  @override
  LoginPageState createState() => LoginPageState();
}

class LoginPageState extends State<LoginPage> {
  final _usernameController = TextEditingController(); // 用户名输入控制器
  final _passwordController = TextEditingController(); // 密码输入控制器
  String _ipaddress = "Loading..."; // IP地址信息
  final AuthService _authService = AuthService(); // 认证服务
  final logger = getLogger(); // 日志记录器

  bool _obscureText = true; // 控制密码可见性

  @override
  void initState() {
    super.initState();
    _getProvinceInfo(); // 初始化时获取省份信息
  }

  @override
  void dispose() {
    _usernameController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  // 切换密码可见性
  void _toggleObscureText() {
    setState(() {
      _obscureText = !_obscureText;
    });
  }

  // 登录方法
  void _login() async {
    if (!mounted) return;

    try {
      logger.d('用户名输入: ${_usernameController.text}');
      logger.d('密码输入: ${_passwordController.text}');

      LoginRequest request = LoginRequest()
        ..username = _usernameController.text
        ..password = _passwordController.text
        ..deviceId = 'web'
        ..timestamp = Timestamp.create()
        ..nonce = Random().nextInt(1000000).toString()
        ..ipaddress = _ipaddress;

      final loginResponse = await _authService.login(request);
      if (!mounted) return;

      logger.i('登录响应: $loginResponse');

      if (loginResponse.statusCode == 0) {
        // 登录成功，导航到主页
        Navigator.pushReplacementNamed(context, '/me');
      } else {
        // 登录失败，显示错误信息
        String errorMessage;
        if (loginResponse.statusCode == authUserLoginFailedCode) {
          errorMessage = authUserLoginFailed;
        } else {
          errorMessage = '登录失败: ${loginResponse.statusMsg}';
        }
        _showErrorSnackBar(errorMessage);
      }
    } catch (e) {
      if (!mounted) return;
      _showErrorSnackBar('登录失败: $e');
    }
  }

  // 获取省份信息
  Future<void> _getProvinceInfo() async {
    String province = await IPLocationUtil.getProvince();
    setState(() {
      _ipaddress = province;
    });
  }

  // 显示错误提示
  void _showErrorSnackBar(String message) {
    Fluttertoast.showToast(
      msg: message,
      toastLength: Toast.LENGTH_SHORT,
      gravity: ToastGravity.BOTTOM,
      timeInSecForIosWeb: 1,
      backgroundColor: Colors.red,
      textColor: Colors.white,
      fontSize: 16.0,
    );
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
              _buildHeader(),
              _buildAvatar(),
              _buildInputFields(),
              _buildLoginButton(),
              _buildRegisterButton(),
            ],
          ),
        ],
      ),
    );
  }

  // 构建页面头部
  Widget _buildHeader() {
    return const Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SizedBox(height: 120),
        Text(
          '欢迎来到Aorb',
          style: TextStyle(
            fontSize: 40,
            fontWeight: FontWeight.bold,
          ),
        ),
        SizedBox(height: 8),
        Text(
          '登录账户解锁更多功能~',
          style: TextStyle(
            fontSize: 20,
            color: Colors.black87,
            fontWeight: FontWeight.bold,
          ),
        ),
        SizedBox(height: 80),
      ],
    );
  }

  // 构建头像
  Widget _buildAvatar() {
    return CircleAvatar(
      radius: 50,
      backgroundColor: Colors.grey[300],
      child: Icon(
        Icons.person,
        size: 50,
        color: Colors.grey[600],
      ),
    );
  }

  // 构建输入字段
  Widget _buildInputFields() {
    return Column(
      children: [
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
      ],
    );
  }

  // 构建登录按钮
  Widget _buildLoginButton() {
    return Column(
      children: [
        const SizedBox(height: 24),
        MaterialButton(
          color: Colors.blue,
          textColor: Colors.white,
          padding: const EdgeInsets.all(16),
          shape: const CircleBorder(),
          onPressed: _login,
          child: const Icon(Icons.arrow_forward),
        ),
      ],
    );
  }

  // 构建注册按钮
  Widget _buildRegisterButton() {
    return TextButton(
      onPressed: () {
        Navigator.of(context).push(
          MaterialPageRoute(
            builder: (context) => const RegisterPage(),
          ),
        );
      },
      child: const Text('没有账户？注册'),
    );
  }
}
