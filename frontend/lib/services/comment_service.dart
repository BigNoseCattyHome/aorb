import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/comment.pbgrpc.dart';
import 'package:grpc/grpc.dart';

class CommentService {
  late final CommentServiceClient _client;
  late final ClientChannel _channel;
  var logger = getLogger();

  CommentService() {
    const host = backendHost;
    const port = backendPort;
    logger.i(
        'Comment Service attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    _client = CommentServiceClient(_channel);
  }

  // ActionCoemment 操作评论（新增/删除）
  Future<ActionCommentResponse> actionComment(
      String username, String pollId, ActionCommentType actionType, String commentText) async {
    logger.i('Comment Service actionComment called');
    final request = ActionCommentRequest()
      ..username = username
      ..pollUuid = pollId
      ..actionType = actionType
      ..commentText = commentText;
    try {
      final response = await _client.actionComment(request);
      return response;
    } catch (e) {
      logger.e('Comment Service actionComment failed', e);
      rethrow;
    }
  }

  // ListComment 获取评论列表
  Future<ListCommentResponse> listComment(String pollId) async {
    logger.i('Comment Service listComment called');
    final request = ListCommentRequest()..pollUuid = pollId;
    try {
      final response = await _client.listComment(request);
      return response;
    } catch (e) {
      logger.e('Comment Service listComment failed', e);
      rethrow;
    }
  }
}
