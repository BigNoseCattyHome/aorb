import 'package:aorb/models/comment.dart';
import 'package:aorb/models/vote.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/mockito.dart';
import 'package:aorb/screens/poll_content.dart';
import 'package:aorb/models/user.dart';
import 'package:aorb/models/poll.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/services/poll_service.dart';

class MockUserService extends Mock implements UserService {}

class MockPollService extends Mock implements PollService {}

void main() {
  late MockUserService mockUserService;
  late MockPollService mockPollService;

  setUp(() {
    mockUserService = MockUserService();
    mockPollService = MockPollService();
  });

  testWidgets('PollContent widget test', (WidgetTester tester) async {
    // 配置mock对象的行为
    when(mockUserService.fetchFollowStatus('1', '2'))
        .thenAnswer((_) async => true); // 1 关注了 2
    when(mockUserService.fetchUserInfo('1', ['nickname', 'avatar'])).thenAnswer(
        (_) async => {'nickname': 'test_user', 'avatar': 'test_avatar'});
    when(mockPollService.fetchPoll('1')).thenAnswer((_) async => Poll(
        id: '1',
        type: 'public',
        sponsor: '1',
        votes: [Vote(choice: '大火腿', userId: '3'),Vote(choice: '火锅', userId: '4')],
        title: '午饭吃什么呀？',
        description: '想了半天没有想出来到底要吃什么，好纠结，真可恶！再不吃饭我真的就要饿死了，下午还要上课，我宝贵的吃饭时间啊啊，来人速速帮我决定一下！',
        options: ['火锅','大火腿'],
        options_rate: [0.3,0.7],
        time: DateTime.parse('2024-01-01 00:00:00'),
        ipaddress: '上海',
        comments: [Comment(advise: '我觉得火锅比较好吃耶，虽然火腿很香，有一点想吃mamamiya了哈哈哈，下次要一起去吗？', choose: '火锅', userid: '3')]));

    // 构建Widget
    await tester.pumpWidget(
      const MaterialApp(
        home: PollContent(
          userID: '1',
          postUserID: '2',
        ),
      ),
    );

    // 确保页面构建完成
    await tester.pump();

    // 验证页面上显示的内容
    expect(find.text('爱吃饭的小袁同学'), findsOneWidget);
    expect(find.text('想了半天没有想出来到底要吃什么，好纠结，真可恶！再不吃饭我真的就要饿死了，下午还要上课，我宝贵的吃饭时间啊啊，来人速速帮我决定一下！'), findsOneWidget);
  });
}
