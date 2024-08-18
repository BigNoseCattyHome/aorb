import 'package:aorb/conf/config.dart';
import 'package:aorb/utils/constant/err.dart';
import 'package:grpc/grpc.dart';
import 'package:aorb/generated/poll.pbgrpc.dart';

class PollService {
  late final PollServiceClient _client;
  late final ClientChannel _channel;
  var logger = getLogger();

  PollService() {
    const host = backendHost;
    const port = backendPort;
    logger.i('PollService attempting to connect to $backendHost:$backendPort');

    _channel = ClientChannel(
      host,
      port: port,
      options: const ChannelOptions(credentials: ChannelCredentials.insecure()),
    );
    _client = PollServiceClient(_channel);
  }

  // 创建投票
  Future<CreatePollResponse> createPoll(CreatePollRequest request) async {
    try {
      final CreatePollResponse response = await _client.createPoll(request);

      if (response.statusCode == 0) {
        logger.i('Successfully created poll ${request.poll.title}');
        return response;
      } else {
        logger.w(
            'Failed to create poll ${request.poll.title}: ${response.statusMsg}');
        throw Exception(
            'Failed to create poll ${request.poll.title}: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('gRPC error occurred while creating poll: ${e.message}');
      throw Exception('Failed to create poll: ${e.message}');
    } catch (e) {
      logger.e('Unexpected error occurred while creating poll: $e');
      throw Exception('Failed to create poll: $e');
    }
  }

  // 查询投票信息
  Future<GetPollResponse> getPoll(GetPollRequest request) async {
    final GetPollResponse response = await _client.getPoll(request);

    if (response.statusCode == 0) {
      logger
          .i('Successfully retrieved information for poll ${request.pollUuid}');
      return response;
    } else {
      logger.w(
          'Failed to retrieve information for poll ${request.pollUuid}: ${response.statusMsg}');
      throw Exception(
          'Failed to retrieve information for poll ${request.pollUuid}: ${response.statusMsg}');
    }
  }

  // feedPoll 推送最新的十条poll
  Future<FeedPollResponse> feedPoll(FeedPollRequest request) async {
    final FeedPollResponse response = await _client.feedPoll(request);

    if (response.statusCode == 0) {
      logger.i('Successfully retrieved feed');
      return response;
    } else {
      logger.w('Failed to retrieve feed: ${response.statusMsg}');
      throw Exception('Failed to retrieve feed: ${response.statusMsg}');
    }
  }

  // GetChoiceWithPollUuidAndUsername 获取用户的投票信息
  Future<GetChoiceWithPollUuidAndUsernameResponse>
      getChoiceWithPollUuidAndUsername(
          GetChoiceWithPollUuidAndUsernameRequest request) async {
    final GetChoiceWithPollUuidAndUsernameResponse response =
        await _client.getChoiceWithPollUuidAndUsername(request);

    if (response.statusCode == 0) {
      return response;
    } else if (response.statusCode == unableToGetChoiceCode) {
      return response;
    } else {
      logger.w(
          'Failed to retrieve choice for poll ${request.pollUuid}: ${response.statusMsg}');
      throw Exception(
          'Failed to retrieve choice for poll ${request.pollUuid}: ${response.statusMsg}');
    }
  }
}
