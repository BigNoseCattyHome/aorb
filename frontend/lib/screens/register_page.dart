import 'package:flutter/material.dart';
import 'package:aorb/conf/config.dart';
import 'package:aorb/services/auth_service.dart';
import 'package:aorb/utils/ip_locator.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

  @override
  RegisterPageState createState() => RegisterPageState();
}

class RegisterPageState extends State<RegisterPage> {
  final _usernameController = TextEditingController(); // 用于控制输入框的值
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  bool _agreeToTerms = false; // 是否同意用户隐私政策条款
  bool _obscureText = true; // 是否隐藏密码
  String _province = 'Loading...'; // 用户IP的归属地
  // 在页面构建的时候初始化AuthService
  final AuthService _authService = AuthService(backendHost, backendPort);
  final logger = getLogger();

  @override
  void initState() {
    super.initState();
    _getProvinceInfo();
  }

  // 调用utils/ip_locator.dart中的getProvince方法获取用户IP的归属地
  Future<void> _getProvinceInfo() async {
    String province = await IPLocationUtil.getProvince();
    setState(() {
      _province = province;
    });
  }

  // 切换密码可见性
  void _toggleObscureText() {
    setState(() {
      _obscureText = !_obscureText;
    });
  }

  // 注册逻辑
  void _register() async {
    // 检查是否同意用户隐私政策条款
    if (!_agreeToTerms) {
      showDialog(
        context: context,
        builder: (BuildContext context) {
          return AlertDialog(
            title: const Text('请同意用户隐私政策条款'),
            content: const Text('您需要同意用户隐私政策条款才能继续注册。'),
            actions: <Widget>[
              TextButton(
                child: const Text('取消'),
                onPressed: () {
                  Navigator.of(context).pop(); // 关闭对话框
                },
              ),
              TextButton(
                child: const Text('同意'),
                onPressed: () {
                  setState(() {
                    _agreeToTerms = true;
                  });
                  Navigator.of(context).pop();
                },
              ),
            ],
          );
        },
      );
      return;
    }

    // 检查两次输入的密码是否一致
    if (_passwordController.text != _confirmPasswordController.text) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('两次输入的密码不一致')),
      );
      return;
    }

    // 调用AuthService的register方法进行注册
    try {
      // ~ 控制台输出后端的主机和端口号
      logger.i('backendHost: $backendHost, backendPort: $backendPort');

      // username和password是必填项，nickname, avatar, ipaddress是可选项
      final registerResponse = await _authService.register(
        _usernameController.text,
        _passwordController.text,
        nickname: _usernameController.text,
        ipaddress: _province,
        avatar: '',
      );

      // ~ 控制台输出注册结果
      logger.i('registerResponse: $registerResponse');
      
      if (registerResponse.success) {
        // 注册成功后的处理逻辑
        Navigator.pop(context);
      } else {
        // 注册失败后的处理逻辑，一般不会走到这里，弹出错误提示
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('注册失败，请检查输入或网络连接。')),
        );
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('注册失败: $e')),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Colors.white,
        shadowColor: Colors.transparent,
        leading: IconButton(
          icon: Icon(
            Icons.arrow_back,
            color: Colors.blue[700],
          ),
          onPressed: () {
            Navigator.pop(context);
          },
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const SizedBox(height: 40),
              const Text(
                '注册',
                style: TextStyle(fontSize: 40, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 8),
              const Text(
                'You and your friends always connected',
                style: TextStyle(color: Colors.grey),
              ),
              const SizedBox(height: 40),
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
              const SizedBox(height: 16),
              TextFormField(
                controller: _confirmPasswordController,
                obscureText: _obscureText,
                decoration: InputDecoration(
                  labelText: '再次输入密码',
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
              CheckboxListTile(
                title: const Text('我同意用户隐私政策条款'),
                value: _agreeToTerms,
                onChanged: (bool? value) {
                  setState(() {
                    _agreeToTerms = value ?? false;
                  });
                },
                controlAffinity: ListTileControlAffinity.leading,
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: _register, // 点击按钮时触发注册的逻辑
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                ),
                child: const Text(
                  '注册',
                ),
              ),
              const SizedBox(height: 16),
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Text('已经有账户了？ '),
                  TextButton(
                    child: const Text('前往登录页面'),
                    onPressed: () {
                      Navigator.pop(context);
                    },
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
