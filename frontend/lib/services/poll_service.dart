import 'package:aorb/conf/config.dart';
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
        logger.i('成功创建投票 ${request.poll.title}');
        return response;
      } else {
        logger.w('创建投票 ${request.poll.title} 失败: ${response.statusMsg}');
        throw Exception('创建投票 ${request..poll.title} 失败: ${response.statusMsg}');
      }
    } on GrpcError catch (e) {
      logger.e('创建投票时发生gRPC错误: ${e.message}');
      throw Exception('创建投票失败: ${e.message}');
    } catch (e) {
      logger.e('创建投票时发生意外错误: $e');
      throw Exception('创建投票失败: $e');
    }
  }

  // 查询投票信息
  Future<GetPollResponse> getPoll(GetPollRequest requset) async {
    final GetPollResponse response = await _client.getPoll(requset);

    if (response.statusCode == 0) {
      logger.i('成功获取投票 ${requset.pollUuid} 的信息');
      return response;
    } else {
      logger.w('获取投票 ${requset.pollUuid} 信息失败: ${response.statusMsg}');
      throw Exception('获取投票 ${requset.pollUuid} 信息失败: ${response.statusMsg}');
    }
  }
}
