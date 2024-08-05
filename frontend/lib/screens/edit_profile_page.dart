import 'package:flutter/material.dart';
import 'package:aorb/generated/user.pbgrpc.dart';
import 'package:image_picker/image_picker.dart';

class EditProfilePage extends StatefulWidget {
  final User user;

  EditProfilePage({required this.user});

  @override
  _EditProfilePageState createState() => _EditProfilePageState();
}

class _EditProfilePageState extends State<EditProfilePage> {
  late User _user;
  final ImagePicker _picker = ImagePicker();

  @override
  void initState() {
    super.initState();
    _user = widget.user;
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('编辑个人资料', style: TextStyle(color: Colors.black, fontWeight: FontWeight.bold)),
        backgroundColor: Colors.white,
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.black),
      ),
      body: ListView(
        children: [
          _buildAvatarSection(),
          _buildDivider(),
          _buildEditItem('昵称', _user.nickname, () => _editNickname()),
          _buildEditItem('用户名', _user.username, () => _editUsername()),
          _buildEditItem('个人简介', _user.bio, () => _editBio()),
          _buildEditItem('性别', _genderToString(_user.gender), () => _editGender()),
          _buildImageItem('背景图片', _user.bgpicMe, () => _editBgPicMe()),
          _buildImageItem('投票背景图片', _user.bgpicPollcard, () => _editBgPicPollcard()),
        ],
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
              backgroundImage: NetworkImage(_user.avatar),
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
                  onPressed: _editAvatar,
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
          Text(value.isEmpty ? '未设置' : value, style: TextStyle(color: Colors.grey[600])),
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
          child: Image.network(
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

  void _editAvatar() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      // TODO: 实现头像上传逻辑
      setState(() {
        _user.avatar = image.path; // 这里应该是上传后的URL
      });
    }
  }

  void _editNickname() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => EditTextPage(title: '昵称', initialValue: _user.nickname)),
    );
    if (result != null) {
      setState(() {
        _user.nickname = result;
      });
    }
  }

  void _editUsername() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => EditTextPage(title: '用户名', initialValue: _user.username)),
    );
    if (result != null) {
      setState(() {
        _user.username = result;
      });
    }
  }

  void _editBio() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => EditTextPage(title: '个人简介', initialValue: _user.bio, maxLines: 5)),
    );
    if (result != null) {
      setState(() {
        _user.bio = result;
      });
    }
  }

  void _editGender() async {
    final result = await Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => EditGenderPage(initialGender: _user.gender)),
    );
    if (result != null) {
      setState(() {
        _user.gender = result;
      });
    }
  }

  void _editBgPicMe() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      // TODO: 实现图片上传逻辑，获取新的URL
      setState(() {
        _user.bgpicMe = image.path; // 这里应该是上传后的URL
      });
    }
  }

  void _editBgPicPollcard() async {
    final XFile? image = await _picker.pickImage(source: ImageSource.gallery);
    if (image != null) {
      // TODO: 实现图片上传逻辑，获取新的URL
      setState(() {
        _user.bgpicPollcard = image.path; // 这里应该是上传后的URL
      });
    }
  }
}

class EditTextPage extends StatefulWidget {
  final String title;
  final String initialValue;
  final int maxLines;

  EditTextPage({required this.title, required this.initialValue, this.maxLines = 1});

  @override
  _EditTextPageState createState() => _EditTextPageState();
}

class _EditTextPageState extends State<EditTextPage> {
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
        title: Text('编辑${widget.title}', style: const TextStyle(color: Colors.black)),
        backgroundColor: Colors.white,
        elevation: 0,
        iconTheme: const IconThemeData(color: Colors.black),
        actions: [
          TextButton(
            onPressed: () {
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
          maxLines: widget.maxLines,
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

  EditGenderPage({required this.initialGender});

  @override
  _EditGenderPageState createState() => _EditGenderPageState();
}

class _EditGenderPageState extends State<EditGenderPage> {
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
        ],
      ),
    );
  }
}