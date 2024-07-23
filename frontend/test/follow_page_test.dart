import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';

import 'package:aorb/screens/follow_page.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/models/user.dart';
import 'package:aorb/models/simple_user.dart';

class MockUserService extends Mock implements UserService {}

void main() {
  late MockUserService mockUserService;

  setUp(() {
    mockUserService = MockUserService();

    // TODO 这里还没有修好！
    // 设置未 stub 方法的默认返回值
    when(mockUserService.fetchUserInfo("")).thenAnswer((_) async => User(
        avatar: '',
        blacklist: [],
        coins: 0,
        coinsRecord: [],
        followed: [],
        follower: [],
        id: '',
        nickname: '',
        ipaddress: '',
        questionsAsk: [],
        questionsAsw: [],
        questionsCollect: []));

    // ! here
    when(mockUserService.fetchFollowList("")).thenAnswer((_) async => []);
    when(mockUserService.fetchFanList("any")).thenAnswer((_) async => []);
  });

  testWidgets('FollowPage widget test', (WidgetTester tester) async {
    // 配置mock对象的行为
    when(mockUserService.fetchUserInfo('1', ['nickname', 'avtar'])).thenAnswer(
        (_) async => User(
            avatar: 'https://s2.loli.net/2024/05/27/2MgJcvLtOVKmAdn.jpg',
            blacklist: [],
            coins: 0,
            coinsRecord: [],
            followed: [],
            follower: [],
            id: '1',
            nickname: '爱吃饭的小袁同学',
            ipaddress: '',
            questionsAsk: [],
            questionsAsw: [],
            questionsCollect: []));
    when(mockUserService.fetchFollowList('1')).thenAnswer((_) async => [
          SimpleUser(
            nickname: "花枝鼠gogo来帮忙",
            avatar: "https://s2.loli.net/2024/05/25/icuYCOP9HB1JbIx.png",
            ipaddress: '2',
          ),
          SimpleUser(
            nickname: "风见澈Siri",
            avatar: "https://s2.loli.net/2024/05/27/QzKM41C3Vs5FeHW.jpg",
            ipaddress: '3',
          )
        ]);
    when(mockUserService.fetchFanList('1')).thenAnswer((_) async => [
          SimpleUser(
            nickname: "Anti Cris",
            avatar: "https://s2.loli.net/2024/05/27/alt3BKPYhzmV4E7.jpg",
            ipaddress: '4',
          ),
        ]);

    // 构建Widget
    await tester.pumpWidget(const MaterialApp(
      home: FollowPage(userId: '1'),
    ));

    // 确保页面构建完成
    await tester.pumpAndSettle();

    // 验证页面上显示的内容
    expect(find.text('花枝鼠gogo来帮忙'), findsOneWidget);
    expect(find.text('风见澈Siri'), findsOneWidget);

    // 找到 TabBar 并模拟点击第二个 Tab
    final tabBar = find.byType(TabBar);
    await tester.tap(tabBar.at(1)); // 假设粉丝页面是第二个 Tab

    // 等待动画完成
    await tester.pumpAndSettle();

    // 检查粉丝页面有没有 "Anti Cris"
    expect(find.text('Anti Cris'), findsOneWidget);
  });
}
