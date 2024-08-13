// home_page.dart
import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/poll.pbgrpc.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/screens/content_publish_page.dart';
import 'package:aorb/screens/login_prompt_page.dart';
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
  late TabController _tabController;
  late List<PollCard> _polls;
  final ScrollController _scrollController = ScrollController();
  bool _isLoading = false;
  bool isLoggedIn = false;
  final logger = getLogger();

  @override
  void initState() {
    super.initState();
    _tabController = widget.tabController;
    _polls = [];
    _fetchPolls();
    _scrollController.addListener(_onScroll);
  }

  @override
  void dispose() {
    _tabController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _fetchPolls({bool isRefresh = false}) async {
    if (_isLoading) return;
    setState(() => _isLoading = true);

    try {
      final response = await PollService()
          .feedPoll(FeedPollRequest()..username = widget.username);

      if (response.statusCode != 0) {
        throw Exception('Failed to fetch polls: ${response.statusMsg}');
      }

      List<PollCard> newPolls = [];
      for (Poll poll in response.pollList) {
        final pollData = await _fetchAdditionalPollData(poll);
        newPolls.add(_createPollCard(poll, pollData));
      }

      setState(() {
        if (isRefresh) {
          _polls = newPolls;
        } else {
          _polls.addAll(newPolls);
        }
        _isLoading = false;
      });
    } catch (e) {
      logger.e('Error fetching polls: $e');
      setState(() => _isLoading = false);
    }
  }

  PollCard _createPollCard(Poll poll, Map<String, dynamic> pollData) {
    final userInfo = pollData['userInfo'] as User;
    final totalVotes = pollData['totalVotes'] as int;
    final percentages = pollData['percentages'] as List<double>;
    final selectedOption = pollData['selectedOption'] as String;

    return PollCard(
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
    );
  }

  void _onScroll() {
    if (_scrollController.position.pixels ==
        _scrollController.position.maxScrollExtent) {
      _fetchPolls();
    }
  }

  Widget _buildPollList() {
    return RefreshIndicator(
      onRefresh: () => _fetchPolls(isRefresh: true),
      child: ListView.builder(
        controller: _scrollController,
        itemCount: _polls.length + 1,
        itemBuilder: (context, index) {
          if (index < _polls.length) {
            return _polls[index];
          } else if (_isLoading) {
            return const Center(child: CircularProgressIndicator());
          } else {
            return const SizedBox.shrink();
          }
        },
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      floatingActionButton: Visibility(
        visible: widget.username != "",
        child: FloatingActionButton(
          onPressed: () {
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
          ),
        ),
      ),
      body: TabBarView(
        controller: _tabController,
        children: [
          _buildPollList(),
          widget.username == "" ? const LoginPromptPage() : _buildPollList(),
        ],
      ),
    );
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
      logger.e('Error fetching additional poll data: $e');
      rethrow;
    }
  }
}
