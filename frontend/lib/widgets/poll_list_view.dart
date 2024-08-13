import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/poll.pb.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/widgets/poll_card.dart';
import 'package:flutter/material.dart';

class PollListView extends StatefulWidget {
  final List<String> pollIds;
  final String currentUsername;
  final String emptyMessage;
  final VoidCallback onRefresh;

  const PollListView({
    Key? key,
    required this.pollIds,
    required this.currentUsername,
    this.emptyMessage = 'No data',
    required this.onRefresh,
  }) : super(key: key);

  @override
  PollListViewState createState() => PollListViewState();
}

class PollListViewState extends State<PollListView> {
  final ScrollController _scrollController = ScrollController();
  final logger = getLogger();

  @override
  void initState() {
    super.initState();
    // 确保在构建完成后滚动到顶部
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.jumpTo(0);
      }
    });
  }

  void refresh() {
    setState(() {
      final reversedPollIds = widget.pollIds.reversed.toList();
      for (int index = 0; index < reversedPollIds.length; index++) {
        final pollId = reversedPollIds[index];
        _fetchPollData(pollId);
      }
    });
  }

  @override
  void dispose() {
    _scrollController.dispose();
    super.dispose();
  }

  Future<Map<String, dynamic>> _fetchPollData(String pollId) async {
    try {
      // 调用PollService获取投票信息
      final pollResponse = await PollService().getPoll(
        GetPollRequest()..pollUuid = pollId,
      );
      final poll = pollResponse.poll;

      // 调用UserService获取用户信息（poll的发起人）
      final userInfoResponse = await UserService().getUserInfo(
        UserRequest()..username = poll.username,
      );
      final userInfo = userInfoResponse.user;

      // 查询用户是否已经投票
      final selectedOptionResponse =
          await PollService().getChoiceWithPollUuidAndUsername(
        GetChoiceWithPollUuidAndUsernameRequest()
          ..pollUuid = pollId
          ..username = widget.currentUsername,
      );
      final selectedOption = selectedOptionResponse.choice;

      // 计算投票百分比
      final totalVotes = poll.optionsCount.reduce((a, b) => a + b);
      final percentages = poll.optionsCount.map((value) {
        return totalVotes > 0 ? (value / totalVotes) * 100 : 0.0;
      }).toList();

      return {
        'poll': poll,
        'userInfo': userInfo,
        'totalVotes': totalVotes,
        'percentages': percentages,
        'selectedOption': selectedOption,
      };
    } catch (e) {
      logger.e('Error fetching poll data: $e');
      rethrow;
    }
  }

  @override
  Widget build(BuildContext context) {
    // 反转pollIds列表，确保最新的在前面
    final reversedPollIds = widget.pollIds.reversed.toList();

    return reversedPollIds.isEmpty
        ? Center(child: Text(widget.emptyMessage))
        : ListView.builder(
            controller: _scrollController,
            itemCount: reversedPollIds.length,
            itemBuilder: (context, index) {
              final pollId = reversedPollIds[index];
              return FutureBuilder<Map<String, dynamic>>(
                future: _fetchPollData(pollId),
                builder: (context, snapshot) {
                  if (snapshot.connectionState == ConnectionState.waiting) {
                    return const Center(child: CircularProgressIndicator());
                  } else if (snapshot.hasError) {
                    return Center(child: Text('Error: ${snapshot.error}'));
                  } else if (snapshot.hasData) {
                    final data = snapshot.data!;
                    final poll = data['poll'] as Poll;
                    final userInfo = data['userInfo'] as User;
                    final totalVotes = data['totalVotes'] as int;
                    final percentages = data['percentages'] as List<double>;
                    final selectedOption = data['selectedOption'] as String;

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
                  } else {
                    return const Center(child: Text('No data'));
                  }
                },
              );
            },
          );
  }
}
