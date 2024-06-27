package config

// 微服务的名字以及端口
const WebServerName = "AorB-Gateway"
const WebServerAddr = ":37000"

const AuthRpcServerName = "AorB-AuthServer"
const AuthRpcServerAddr = ":37001"

const UserRpcServerName = "AorB-UserServer"
const UserRpcServerAddr = ":37002"

const CommentRpcServerName = "AorB-CommentServer"
const CommentRpcServerAddr = ":37003"

const VoteServerName = "AorB-VoteServer"
const VoteServerAddr = ":37004"

const PollRpcServerName = "AorB-PollServer"
const PollRpcServerAddr = ":37005"

const RecommendRpcServerName = "AorB-RecommendServer"
const RecommendRpcServerAddr = ":37006"

const Metrics = ":37099"
const PollProcessorRpcServiceName = "AorB-PollProcessorServer"

const Event = "AorB-Recommend"
const MsgConsumer = "AorB-MgsConsumer"
const BloomRedisChannel = "AorB-Bloom"
