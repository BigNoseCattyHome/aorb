// home_page.dart
import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/poll.pbgrpc.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/screens/content_publish_page.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/services/user_service.dart';
import 'package:flutter/material.dart';
import 'package:aorb/widgets/poll_card.dart';
import 'dart:async';

class HomePage extends StatefulWidget {
  final TabController tabController;
  final String username;

  const HomePage(
      {super.key, required this.tabController, required this.username});

  @override
  HomePageState createState() => HomePageState();
}

class HomePageState extends State<HomePage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController; // 顶部导航栏控制器
  late Future<List<PollCard>> _futurePolls;
  final logger = getLogger();
  @override
  void initState() {
    super.initState();
    _tabController = widget.tabController; // 初始化顶部导航栏控制器
    _futurePolls = _fetchPolls(); // 初始化时调用 _fetchQuestions 获取数据
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      // 顶部栏，由main_page接管

      // 发布按钮
      floatingActionButton: FloatingActionButton(
          onPressed: () {
            // 导航到发布界面
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => const ContentPublishPage(),
              ),
            );
          },
          backgroundColor: Colors.blue[700],
          child: const Icon(
            Icons.add,
            color: Colors.white,
            size: 30,
          )),

      // 中间的投票卡片
      body: TabBarView(controller: _tabController, children: [
        FutureBuilder<List<PollCard>>(
          future: _futurePolls,
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return const Center(child: CircularProgressIndicator());
            } else if (snapshot.hasError) {
              return Center(child: Text('Error: ${snapshot.error}'));
            } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
              return const Center(child: Text('No questions available.'));
            } else {
              return ListView.builder(
                itemCount: snapshot.data!.length,
                itemBuilder: (context, index) {
                  return snapshot.data![index];
                },
              );
            }
          },
        ),
        FutureBuilder<List<PollCard>>(
          future: _futurePolls,
          builder: (context, snapshot) {
            if (snapshot.connectionState == ConnectionState.waiting) {
              return const Center(child: CircularProgressIndicator());
            } else if (snapshot.hasError) {
              return Center(child: Text('Error: ${snapshot.error}'));
            } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
              return const Center(child: Text('No questions available.'));
            } else {
              return ListView.builder(
                itemCount: snapshot.data!.length,
                itemBuilder: (context, index) {
                  return snapshot.data![index];
                },
              );
            }
          },
        ),
      ]),
    );
  }

  Future<List<PollCard>> _fetchPolls() async {
    try {
      final response = await PollService()
          .feedPoll(FeedPollRequest()..username = widget.username);

      if (response.statusCode != 0) {
        throw Exception('Failed to fetch polls: ${response.statusMsg}');
      }

      List<PollCard> pollCards = [];

      for (Poll poll in response.pollList) {
        final pollData = await _fetchAdditionalPollData(poll);
        final userInfo = pollData['userInfo'] as User;
        final totalVotes = pollData['totalVotes'] as int;
        final percentages = pollData['percentages'] as List<double>;
        final selectedOption = pollData['selectedOption'] as String;

        pollCards.add(PollCard(
          pollId: poll.pollUuid,
          title: poll.title,
          content: poll.content,
          options: poll.options,
          voteCount: totalVotes,
          time: poll.createAt,
          username: userInfo.username,
          avatar: userInfo.avatar,
          nickname: userInfo.nickname,
          userId: userInfo.id,
          backgroundImage: userInfo.bgpicPollcard,
          votePercentage: percentages,
          selectedOption: selectedOption,
        ));
      }

      // 如果需要，可以在这里保存 nextTime 以便后续使用
      // final nextTime = response.nextTime;

      return pollCards;
    } catch (e) {
      logger.e('Error fetching polls: $e');
      rethrow;
    }
  }

  Future<Map<String, dynamic>> _fetchAdditionalPollData(Poll poll) async {
    try {
      final userInfoResponse = await UserService().getUserInfo(
        UserRequest()..username = poll.username,
      );
      final userInfo = userInfoResponse.user;

      final selectedOptionResponse =
          await PollService().getChoiceWithPollUuidAndUsername(
        GetChoiceWithPollUuidAndUsernameRequest()
          ..pollUuid = poll.pollUuid
          ..username = widget.username,
      );
      final selectedOption = selectedOptionResponse.choice;

      final totalVotes = poll.optionsCount.reduce((a, b) => a + b);
      final percentages = poll.optionsCount.map((value) {
        return totalVotes > 0 ? (value / totalVotes) * 100 : 0.0;
      }).toList();

      return {
        'userInfo': userInfo,
        'totalVotes': totalVotes,
        'percentages': percentages,
        'selectedOption': selectedOption,
      };
    } catch (e) {
      print('Error fetching additional poll data: $e');
      rethrow;
    }
  }
}
