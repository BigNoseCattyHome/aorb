import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/message.pb.dart';
import 'package:aorb/screens/poll_detail_page.dart';
import 'package:aorb/screens/user_profile_page.dart';
import 'package:aorb/services/message_service.dart';
import 'package:flutter/material.dart';
import 'package:shared_preferences/shared_preferences.dart';

class MessagesPage extends StatefulWidget {
  final TabController tabController;

  const MessagesPage({super.key, required this.tabController});

  @override
  MessagesPageState createState() => MessagesPageState();
}

class MessagesPageState extends State<MessagesPage> {
  late TabController _tabController;
  String username = '';
  String userId = '';
  UserMessageResponse messages = UserMessageResponse();
  bool isLoading = true; // 是否正在加载数据
  final logger = getLogger(); // 日志记录器

  @override
  void initState() {
    super.initState();
    _tabController = widget.tabController;
    _loadData();
  }

  void _loadData() async {
    setState(() => isLoading = true);

    final prefs = await SharedPreferences.getInstance();
    username = prefs.getString('username') ?? '';
    userId = prefs.getString('userId') ?? '';

    try {
      UserMessageResponse messagesFetch =
          await MessageService().getUserMessage(username);
      setState(() {
        messages = messagesFetch;
        isLoading = false;
      });
    } catch (e) {
      logger.e('Error loading messages: $e');
      setState(() => isLoading = false);
    }
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: TabBarView(
        controller: _tabController,
        children: [_buildNoticeList(), _buildMessageList()],
      ),
    );
  }

  Widget _buildNoticeList() {
    if (isLoading) {
      return const Center(child: CircularProgressIndicator());
    }

    return RefreshIndicator(
      onRefresh: () async => _loadData(), // onRefresh 在下拉刷新时调用
      child: ListView.builder(
        itemCount: getTotalLength([
          messages.messagesCommentReplyList,
          messages.messagesFollowList,
          messages.messagesVoteList,
        ]),
        itemBuilder: (context, index) {
          // 根据index判断是那个类型的消息
          final message = _getMessageAtIndex(index);
          return _buildMessageTile(message);
        },
      ),
    );
  }

  Widget _buildMessageTile(dynamic message) {
    return ListTile(
      leading: _getLeadingIcon(message),
      title: Text(message.title ?? ''),
      subtitle: Text(message.content ?? ''),
      trailing: message.isRead
          ? null
          : const Icon(Icons.circle, color: Colors.red, size: 10),
      onTap: () => _handleMessageTap(message),
    );
  }

  Icon _getLeadingIcon(dynamic message) {
    if (message is CommentReplyMessage) {
      return const Icon(Icons.comment);
    } else if (message is FollowMessage) {
      return const Icon(Icons.person_add);
    } else if (message is VoteMessage) {
      return const Icon(Icons.how_to_vote);
    }
    return const Icon(Icons.message);
  }

  void _handleMessageTap(dynamic message) async {
    // 标记消息为已读，await 不会阻塞后续操作
    await MessageService()
        .markMessageStatus(message.message_id, MessageStatus.READ);
    if (!mounted) return; // 防止页面已经被销毁

    // 根据消息类型跳转到不同的页面
    if (message is CommentReplyMessage) {
      // 跳转到问题页面
      Navigator.push(
          context,
          MaterialPageRoute(
              builder: (context) => PollDetailPage(
                    pollId: message.pollId,
                    username: username,
                    userId: userId,
                  )));
    } else if (message is FollowMessage) {
      // 跳转到用户详情页面
      Navigator.push(
          context,
          MaterialPageRoute(
              builder: (context) =>
                  UserProfilePage(username: message.usernameFollower)));
    } else if (message is VoteMessage) {
      // 跳转到问题页面
      Navigator.push(
          context,
          MaterialPageRoute(
              builder: (context) => PollDetailPage(
                    pollId: message.pollId,
                    username: username,
                    userId: userId,
                  )));
    }
  }

  // 根据索引获取消息
  dynamic _getMessageAtIndex(int index) {
    int commentReplyCount = messages.messagesCommentReplyList.length;
    int followCount = messages.messagesFollowList.length;

    if (index < commentReplyCount) {
      return messages.messagesCommentReplyList[index];
    } else if (index < commentReplyCount + followCount) {
      return messages.messagesFollowList[index - commentReplyCount];
    } else {
      return messages.messagesVoteList[index - commentReplyCount - followCount];
    }
  }

  Widget _buildMessageList() {
    // TODO: 实现私信列表
    return const Center(child: Text('私信功能尚未实现'));
  }

  int getTotalLength(List<List<dynamic>> lists) {
    return lists.fold(0, (sum, list) => sum + list.length);
  }
}
