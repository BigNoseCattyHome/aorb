import 'package:aorb/generated/user.pb.dart';
import 'package:flutter/material.dart';
import 'package:aorb/conf/config.dart';
import 'package:aorb/services/auth_service.dart';
import 'package:aorb/utils/ip_locator.dart';
import 'package:grpc/grpc.dart';

class RegisterPage extends StatefulWidget {
  const RegisterPage({super.key});

  @override
  RegisterPageState createState() => RegisterPageState();
}

class RegisterPageState extends State<RegisterPage> {
  final _usernameController = TextEditingController();
  final _passwordController = TextEditingController();
  Gender _selectedGender = Gender.UNKNOWN;
  final _confirmPasswordController = TextEditingController();
  bool _agreeToTerms = false;
  bool _obscureText = true;
  String _province = 'Loading...';
  final AuthService _authService = AuthService();
  final logger = getLogger();

  @override
  void initState() {
    super.initState();
    _getProvinceInfo();
  }

  Future<void> _getProvinceInfo() async {
    String province = await IPLocationUtil.getProvince();
    setState(() {
      _province = province;
    });
  }

  void _toggleObscureText() {
    setState(() {
      _obscureText = !_obscureText;
    });
  }

  void _register() async {
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
                  Navigator.of(context).pop();
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

    if (_passwordController.text != _confirmPasswordController.text) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('两次输入的密码不一致')),
      );
      return;
    }

    try {
      logger.d('Attempting to connect to backend at $backendHost:$backendPort');
      logger.d('Register request parameters:');
      logger.d('Username: ${_usernameController.text}');
      logger.d('Password: [REDACTED]');
      logger.d('Nickname: ${_usernameController.text}');
      logger.d('IP Address: $_province');
      logger.d('Avatar: [Empty]');
      final startTime = DateTime.now();
      logger.d('Starting register call at $startTime');

      final registerResponse = await _authService.register(
        _usernameController.text,
        _passwordController.text,
        _selectedGender,
        nickname: _usernameController.text,
        ipaddress: _province,
        avatar: '',
      );

      logger.i('registerResponse: $registerResponse');
      final endTime = DateTime.now();
      final duration = endTime.difference(startTime);
      logger.i(
          'Register call completed at $endTime (Duration: ${duration.inMilliseconds}ms)');
      logger.i('Register response:');
      logger.i('statusCode: ${registerResponse.statusCode}');
      logger.i('statusMsg: ${registerResponse.statusMsg}');

      if (registerResponse.statusCode == 0) {
        Navigator.pop(context);
      } else {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('注册失败，请检查输入或网络连接。')),
        );
      }
    } catch (e) {
      logger.e('Exception occurred during registration');
      logger.e('Error type: ${e.runtimeType}');

      if (e is GrpcError) {
        logger.e('gRPC error code: ${e.code}');
        logger.e('gRPC error details: ${e.details}');
        logger.e('gRPC error trailers: ${e.trailers}');
        logger.e('gRPC error message: ${e.message}');
      }
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
                decoration: InputDecoration(
                  labelText: '用户名',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                  prefixIcon: const Icon(Icons.person),
                ),
              ),
              const SizedBox(height: 16),
              DropdownButtonFormField<Gender>(
                value: _selectedGender,
                onChanged: (Gender? newValue) {
                  setState(() {
                    _selectedGender = newValue ?? Gender.UNKNOWN;
                  });
                },
                items:
                    Gender.values.map<DropdownMenuItem<Gender>>((Gender value) {
                  return DropdownMenuItem<Gender>(
                    value: value,
                    child: Text(_getGenderString(value)),
                  );
                }).toList(),
                decoration: InputDecoration(
                  labelText: '性别',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                  prefixIcon: const Icon(Icons.wc),
                ),
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _passwordController,
                obscureText: _obscureText,
                decoration: InputDecoration(
                  labelText: '密码',
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                  prefixIcon: const Icon(Icons.lock),
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
                  border: OutlineInputBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                  prefixIcon: const Icon(Icons.lock),
                  suffixIcon: IconButton(
                    icon: Icon(
                      _obscureText ? Icons.visibility : Icons.visibility_off,
                    ),
                    onPressed: _toggleObscureText,
                  ),
                ),
              ),
              const SizedBox(height: 16),
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
                onPressed: _register,
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(10),
                  ),
                ),
                child: const Text(
                  '注册',
                  style: TextStyle(fontSize: 18),
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

  String _getGenderString(Gender gender) {
    switch (gender) {
      case Gender.MALE:
        return '男性';
      case Gender.FEMALE:
        return '女性';
      case Gender.OTHER:
        return '其他';
      case Gender.UNKNOWN:
      default:
        return '未知';
    }
  }
}
