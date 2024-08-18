import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/message.pb.dart';
import 'package:aorb/screens/login_prompt_page.dart';
import 'package:aorb/services/message_service.dart';
import 'package:aorb/utils/auth_provider.dart';
import 'package:aorb/widgets/message_comment_reply.dart';
import 'package:aorb/widgets/message_followed.dart';
import 'package:aorb/widgets/message_vote.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:shared_preferences/shared_preferences.dart';

class MessagesPage extends StatefulWidget {
  final TabController tabController;
  final String username;
  const MessagesPage(
      {super.key, required this.tabController, required this.username});

  @override
  MessagesPageState createState() => MessagesPageState();
}

class MessagesPageState extends State<MessagesPage> {
  late TabController _tabController;
  String username = '';
  String userId = '';
  GetUserMessageResponse messages = GetUserMessageResponse();
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
      GetUserMessageResponse messagesFetch =
          await MessageService().getUserMessage(username);
      logger.i('messagesFetch: $messagesFetch');
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
    final authProvider = Provider.of<AuthProvider>(context);

    return Scaffold(
      body: TabBarView(
        controller: _tabController,
        children: [
          authProvider.isLoggedIn
              ? _buildNoticeList()
              : const LoginPromptPage(),
          authProvider.isLoggedIn
              ? _buildMessageList()
              : const LoginPromptPage(),
        ],
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
          messages.commentReplyMessages,
          messages.followMessages,
          messages.voteMessages,
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
    if (message is CommentReplyMessage) {
      return MessageCommentReply(
        time: message.timestamp,
        pollId: message.pollUuid,
        replyText: message.content,
        username: message.username,
        onRead: () {
          _sendReadMessage(message.messageUuid);
        },
      );
    } else if (message is FollowMessage) {
      return MessageFollowed(
        username: message.usernameFollower,
        time: message.timestamp,
        onRead: () {
          _sendReadMessage(message.messageUuid);
        },
      );
    } else if (message is VoteMessage) {
      return MessageVote(
        username: message.voteUsername,
        time: message.timestamp,
        choice: message.choice,
        pollId: message.pollUuid,
        onRead: () {
          _sendReadMessage(message.messageUuid);
        },
      );
    } else {
      return const ListTile(
        title: Text('Unknown message type'),
        subtitle: Text('This message type is not supported yet'),
      );
    }
  }

  // 发送已读消息
  void _sendReadMessage(dynamic message) async {
    await MessageService().markMessageStatus(
        message.message_id, MessageStatus.MESSAGE_STATUS_READ);
    logger.d('Message marked as read: ${message.message_id}');
  }

  // 根据索引获取消息
  dynamic _getMessageAtIndex(int index) {
    int commentReplyCount = messages.commentReplyMessages.length;
    int followCount = messages.followMessages.length;

    if (index < commentReplyCount) {
      return messages.commentReplyMessages[index];
    } else if (index < commentReplyCount + followCount) {
      return messages.followMessages[index - commentReplyCount];
    } else {
      return messages.voteMessages[index - commentReplyCount - followCount];
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
