package strings

// Exchange name
const (
	EventExchange   = "event"
	MessageExchange = "message_exchange"
	AuditExchange   = "audit_exchange"
)

// Queue name
const (
	MessageCommon = "message_common"
	MessageES     = "message_es"
	AuditPicker   = "audit_picker"
)

// Routing key
const (
	FavoriteActionEvent = "poll.favorite.action"
	PollGetEvent        = "poll.get.action"
	PollCommentEvent    = "poll.comment.action"
	PollPublishEvent    = "poll.publish.action"

	MessageActionEvent    = "message.common"
	MessageGptActionEvent = "message.gpt"
	AuditPublishEvent     = "audit"
)

// Action Type
const (
	FavoriteIdActionLog = 1 // 用户点赞相关操作
	FollowIdActionLog   = 2 // 用户关注相关操作
)

// Action Name
const (
	FavoriteNameActionLog    = "favorite.action" // 用户点赞操作名称
	FavoriteUpActionSubLog   = "up"
	FavoriteDownActionSubLog = "down"

	FollowNameActionLog    = "follow.action" // 用户关注操作名称
	FollowUpActionSubLog   = "up"
	FollowDownActionSubLog = "down"
)

// Action Service Name
const (
	FavoriteServiceName = "FavoriteService"
	FollowServiceName   = "FollowService"
)
