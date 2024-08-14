import 'package:aorb/services/auth_service.dart';
import 'package:flutter/material.dart';

class AuthProvider extends ChangeNotifier {
  bool _isLoggedIn = false;

  AuthProvider() {
    _initializeLoginStatus();
  }

  Future<void> _initializeLoginStatus() async {
    _isLoggedIn =
        await AuthService().checkLoginStatus(); // 通过检查 authTokens 是否存在来判断是否登录
    notifyListeners();
  }

  bool get isLoggedIn => _isLoggedIn;

  Future<void> login() async {
    _isLoggedIn = true;
    notifyListeners();
  }

  Future<void> logout() async {
    _isLoggedIn = false;
    notifyListeners();
  }
}
