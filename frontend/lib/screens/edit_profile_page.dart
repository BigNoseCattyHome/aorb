import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:image_picker/image_picker.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/utils/image_upload.dart';
import 'dart:io';

import 'package:shared_preferences/shared_preferences.dart';

class EditProfilePage extends StatefulWidget {
  final User user;
  final Function(User) onUserUpdated;
  const EditProfilePage(
      {Key? key, required this.user, required this.onUserUpdated})
      : super(key: key);

  @override
  EditProfilePageState createState() => EditProfilePageState();
}

class EditProfilePageState extends State<EditProfilePage> {
  late User _user;
  final ImagePicker _picker = ImagePicker();

  @override
  void initState() {
    super.initState();
    _user = widget.user;
  }

  void _updateUser() {
    widget.onUserUpdated(_user);
  }

  @override
  Widget build(BuildContext context) {
    return WillPopScope(
      onWillPop: () async {
        Navigator.pop(context, _user);
        return false;
      },
      child: Scaffold(
        appBar: AppBar(
          title: const Text('编辑个人资料',
              style:
                  TextStyle(color: Colors.black, fontWeight: FontWeight.bold)),
          backgroundColor: Colors.white,
          elevation: 0,
          iconTheme: const IconThemeData(color: Colors.black),
        ),
        body: ListView(
          children: [
            _buildAvatarSection(),
            _buildDivider(),
            _buildEditItem(
                '昵称',
                _user.nickname,
                () => _editField('昵称', _user.nickname, "nickname",
                    (value) => _user.nickname = value)),
            _buildEditItem(
                '用户名',
                _user.username,
                () => _editField('用户名', _user.username, "username",
                    (value) => _user.username = value)),
            _buildEditItem(
                '个人简介',
                _user.bio,
                () => _editField(
                    '个人简介', _user.bio, "bio", (value) => _user.bio = value)),
            _buildEditItem(
                '性别', _genderToString(_user.gender), () => _editGender()),
            _buildImageItem('背景图片', _user.bgpicMe, () => _editImage('bgpicMe')),
            _buildImageItem('投票背景图片', _user.bgpicPollcard,
                () => _editImage('bgpicPollcard')),
          ],
        ),
      ),
    );
  }

  Widget _buildAvatarSection() {
    return Container(
      padding: const EdgeInsets.symmetric(vertical: 20),
      child: Center(
        child: Stack(
          children: [
            CircleAvatar(
              radius: 60,
              backgroundImage: _user.avatar.startsWith('http')
                  ? NetworkImage(_user.avatar)
                  : FileImage(File(_user.avatar)) as ImageProvider,
            ),
            Positioned(
              bottom: 0,
              right: 0,
              child: Container(
                decoration: const BoxDecoration(
                  color: Colors.blue,
                  shape: BoxShape.circle,
                ),
                child: IconButton(
                  icon: const Icon(Icons.edit, color: Colors.white),
                  onPressed: () => _editImage('avatar'),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDivider() {
    return Divider(height: 1, thickness: 1, color: Colors.grey[300]);
  }

  Widget _buildEditItem(String title, String value, VoidCallback onTap) {
    return ListTile(
      title: Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
      trailing: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Text(value.isEmpty ? '未设置' : value,
              style: TextStyle(color: Colors.grey[600])),
          const Icon(Icons.chevron_right, color: Colors.grey),
        ],
      ),
      onTap: onTap,
    );
  }

  Widget _buildImageItem(String title, String imageUrl, VoidCallback onTap) {
    return ListTile(
      title: Text(title, style: const TextStyle(fontWeight: FontWeight.bold)),
      trailing: GestureDetector(
        onTap: onTap,
        child: ClipRRect(
          borderRadius: BorderRadius.circular(8),
          child: imageUrl.startsWith('http')
              ? Image.network(
                  imageUrl,
                  width: 60,
                  height: 60,
                  fit: BoxFit.cover,
                  errorBuilder: (context, error, stackTrace) {
                    return Container(
                      width: 60,
                      height: 60,
                      color: Colors.grey[300],
                      child: const Icon(Icons.image, color: Colors.grey),
                    );
                  },
                )
              : Image.file(
                  File(imageUrl),
                  width: 60,
                  height: 60,
                  fit: BoxFit.cover,
                ),
        ),
      ),
    );
  }

  String _genderToString(Gender gender) {
    switch (gender) {
      case Gender.MALE:
        return '男';
      case Gender.FEMALE:
        return '女';
      case Gender.OTHER:
        return '其他';
      default:
        return '未知';
    }
  }

  void _editField(String title, String initialValue, String type,
      Function(String) onSave) async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(
          builder: (context) => EditTextPage(
              title: title,
              initialValue: initialValue,
              type: type,
              userId: _user.id)),
    );
    if (result != null) {
      setState(() {
        onSave(result);
        _updateUser();
      });

      // 更新shared_preferences
      SharedPreferences prefs = await SharedPreferences.getInstance();
      prefs.setString(type, initialValue);
    }
  }

  void _editGender() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(
          builder: (context) =>
              EditGenderPage(initialGender: _user.gender, userId: _user.id)),
    );
    if (result != null) {
      setState(() {
        _user.gender = result;
        _updateUser();
      });
    }
  }

  void _editImage(String field) async {
    // 选择图片
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);

    // 如果图像大小超过 5MB，提示重新上传
    if (image != null) {
      final File file = File(image.path);
      final size = file.lengthSync();
      if (size > 5 * 1024 * 1024) {
        Fluttertoast.showToast(
          msg: '图片大小不能超过 5MB',
          toastLength: Toast.LENGTH_SHORT,
          gravity: ToastGravity.BOTTOM,
          timeInSecForIosWeb: 1,
          backgroundColor: Colors.blue,
          textColor: Colors.white,
          fontSize: 16.0,
        );
        return;
      }
    }

    if (image != null) {
      // 将图片上传到图床
      final imageUrl = await ImageUploadService()
          .uploadImage(File(image.path), '${field}_${_user.username}');

      // 更新shared_preferences
      SharedPreferences prefs = await SharedPreferences.getInstance();
      prefs.setString(field, imageUrl);

      // 更新用户信息
      setState(() {
        switch (field) {
          case 'avatar':
            _user.avatar = imageUrl;
            break;
          case 'bgpicMe':
            _user.bgpicMe = imageUrl;
            break;
          case 'bgpicPollcard':
            _user.bgpicPollcard = imageUrl;
            break;
        }
        if (field == 'avatar') {}
        _updateUser();
      });
      // 构建UpdateUserRequest
      var request = UpdateUserRequest();
      request.userId = _user.id;
      switch (field) {
        case 'avatar':
          request.avatar = imageUrl;
          break;
        case 'bgpicMe':
          request.bgpicMe = imageUrl;
          break;
        case 'bgpicPollcard':
          request.bgpicPollcard = imageUrl;
          break;
      }
      // 发送gRPC请求更新用户信息
      UserService().updateUser(request);
    }
  }
}

class EditTextPage extends StatefulWidget {
  final String title;
  final String initialValue;
  final String type;
  final String userId;

  const EditTextPage({
    Key? key,
    required this.title,
    required this.initialValue,
    required this.type,
    required this.userId,
  }) : super(key: key);

  @override
  EditTextPageState createState() => EditTextPageState();
}

class EditTextPageState extends State<EditTextPage> {
  late TextEditingController _controller;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(text: widget.initialValue);
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('编辑${widget.title}',
            style: const TextStyle(color: Colors.black)),
        backgroundColor: Colors.white,
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.black),
        actions: [
          TextButton(
            onPressed: () {
              var request = UpdateUserRequest();
              request.userId = widget.userId;
              switch (widget.type) {
                case 'nickname':
                  request.nickname = _controller.text;
                  break;
                case 'bio':
                  request.bio = _controller.text;
                  break;
                case 'username':
                  request.username = _controller.text;
                  break;
              }
              UserService().updateUser(request);
              Navigator.pop(context, _controller.text);
            },
            child: const Text('保存', style: TextStyle(color: Colors.blue)),
          ),
        ],
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: TextField(
          controller: _controller,
          decoration: InputDecoration(
            border: const OutlineInputBorder(),
            labelText: widget.title,
          ),
        ),
      ),
    );
  }
}

class EditGenderPage extends StatefulWidget {
  final Gender initialGender;
  final String userId;

  const EditGenderPage(
      {Key? key, required this.initialGender, required this.userId})
      : super(key: key);

  @override
  EditGenderPageState createState() => EditGenderPageState();
}

class EditGenderPageState extends State<EditGenderPage> {
  late Gender _selectedGender;

  @override
  void initState() {
    super.initState();
    _selectedGender = widget.initialGender;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('编辑性别', style: TextStyle(color: Colors.black)),
        backgroundColor: Colors.white,
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.black),
        actions: [
          TextButton(
            onPressed: () {
              var request = UpdateUserRequest()
                ..userId = widget.userId
                ..gender = _selectedGender;
              UserService().updateUser(request);
              Navigator.pop(context, _selectedGender);
            },
            child: const Text('保存', style: TextStyle(color: Colors.blue)),
          ),
        ],
      ),
      body: Column(
        children: [
          RadioListTile<Gender>(
            title: const Text('男'),
            value: Gender.MALE,
            groupValue: _selectedGender,
            onChanged: (Gender? value) {
              setState(() {
                _selectedGender = value!;
              });
            },
          ),
          RadioListTile<Gender>(
            title: const Text('女'),
            value: Gender.FEMALE,
            groupValue: _selectedGender,
            onChanged: (Gender? value) {
              setState(() {
                _selectedGender = value!;
              });
            },
          ),
          RadioListTile<Gender>(
            title: const Text('其他'),
            value: Gender.OTHER,
            groupValue: _selectedGender,
            onChanged: (Gender? value) {
              setState(() {
                _selectedGender = value!;
              });
            },
          ),
          RadioListTile<Gender>(
            title: const Text('未知'),
            value: Gender.UNKNOWN,
            groupValue: _selectedGender,
            onChanged: (Gender? value) {
              setState(() {
                _selectedGender = value!;
              });
            },
          ),
        ],
      ),
    );
  }
}
