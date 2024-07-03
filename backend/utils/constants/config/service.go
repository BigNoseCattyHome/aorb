package config

// 微服务的名字以及端口
const WebServerName = "AorB-WebGateway"
const WebServerAddr = ":37000"

const AuthRpcServerName = "AorB-AuthService"
const AuthRpcServerAddr = ":37001"

const UserRpcServerName = "AorB-UserService"
const UserRpcServerAddr = ":37002"

const CommentRpcServerName = "AorB-CommentService"
const CommentRpcServerAddr = ":37003"

const VoteRpcServerName = "AorB-VoteService"
const VoteRpcServerAddr = ":37004"

const PollRpcServerName = "AorB-PollService"
const PollRpcServerAddr = ":37005"

const RecommendRpcServerName = "AorB-RecommendService"
const RecommendRpcServerAddr = ":37006"

const Metrics = ":37099"
const PollProcessorRpcServiceName = "GuGoTik-PollProcessorService"

const Event = "GuGoTik-Recommend"
const MsgConsumer = "GuGoTik-MgsConsumer"
const BloomRedisChannel = "GuGoTik-Bloom"
