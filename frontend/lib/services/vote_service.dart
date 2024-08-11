import 'package:aorb/conf/config.dart';
import 'package:aorb/generated/vote.pbgrpc.dart';
import 'package:grpc/grpc.dart';

class VoteService {
  late final VoteServiceClient _client;
  late final ClientChannel _channel;
  var logger = getLogger();

  VoteService() {
    const host = backendHost;
    const port = backendPort;
    logger.i('Vote Service attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    _client = VoteServiceClient(_channel);
  }

  // CreateVote 创建投票
  Future<CreateVoteResponse> createVote(
      String pollId, String username, String choice) async {
    logger.i('Vote Service createVote called');
    final request = CreateVoteRequest()
      ..pollUuid = pollId
      ..username = username
      ..choice = choice;
    try {
      final response = await _client.createVote(request);
      return response;
    } catch (e) {
      logger.e('Vote Service createVote failed', e);
      rethrow;
    }
  }
}
