import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/message.pbgrpc.dart';
import 'package:grpc/grpc.dart';

class MessageService {
  late final MessageServiceClient _client;
  late final ClientChannel _channel;
  var logger = getLogger();

  MessageService() {
    const host = backendHost;
    const port = backendPort;
    logger.i(
        'Message Service attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    _client = MessageServiceClient(_channel);
  }

  // GetUserMessage 获取用户的未读信息
  Future<UserMessageResponse> getUserMessage(String username) async {
    logger.i('Message Service getUserMessage called');
    final request = UserMessageRequest()..username = username;
    try {
      final response = await _client.getUserMessage(request);
      return response;
    } catch (e) {
      logger.e('Message Service getUserMessage failed', e);
      rethrow;
    }
  }

  // MarkMessageStatus 标记消息状态
  Future<MarkMessageStatusResponse> markMessageStatus(
      String messageId, MessageStatus status) async {
    logger.i('Message Service markMessageStatus called');
    final request = MarkMessageStatusRequest()
      ..messageId = messageId
      ..status = status;
    try {
      final response = await _client.markMessageStatus(request);
      return response;
    } catch (e) {
      logger.e('Message Service markMessageStatus failed', e);
      rethrow;
    }
  }
}
