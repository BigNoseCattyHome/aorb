import 'package:aorb/generated/poll.pb.dart';
import 'package:aorb/generated/user.pb.dart';
import 'package:aorb/services/poll_service.dart';
import 'package:aorb/services/user_service.dart';
import 'package:aorb/widgets/poll_card.dart';
import 'package:flutter/material.dart';

class PollListView extends StatelessWidget {
  final List<String> pollIds;
  final String emptyMessage;

  const PollListView({
    Key? key,
    required this.pollIds,
    this.emptyMessage = 'No data',
  }) : super(key: key);

  Future<Map<String, dynamic>> _fetchPollData(String pollId) async {
    // 调用PollService获取投票信息
    final pollResponse = await PollService().getPoll(
      GetPollRequest()..pollUuid = pollId,
    );
    final poll = pollResponse.poll;

    // 调用UserService获取用户信息
    final userInfoResponse = await UserService().getUserInfo(
      UserRequest()..username = poll.username,
    );
    final userInfo = userInfoResponse.user;

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
    };
  }

  @override
  Widget build(BuildContext context) {
    if (pollIds.isEmpty) {
      return Center(child: Text(emptyMessage));
    }

    return ListView.builder(
      itemCount: pollIds.length,
      itemBuilder: (context, index) {
        final pollId = pollIds[index];
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
