import 'package:aorb/generated/poll.pb.dart';
import 'package:aorb/services/poll_service.dart';

void main() {
  PollService()
      .getPoll(
          GetPollRequest()..pollUuid = '8b46ced6-d760-4494-b632-861124ad2d1f')
      .then((resp) => print(resp.poll.content));
}
