package strings

// Bad Request
const (
	GateWayErrorCode       = 40001
	GateWayError           = "AorB Gateway 暂时不能处理您的请求，请稍后重试！"
	GateWayParamsErrorCode = 40002
	GateWayParamsError     = "AorB Gateway 无法响应您的请求，请重启 APP 或稍后再试!"
)

// Server Inner Error
const (
	AuthServiceInnerErrorCode         = 50001
	AuthServiceInnerError             = "登录服务出现内部错误，请稍后重试！"
	UnableToQueryUserErrorCode        = 50002
	UnableToQueryUserError            = "无法查询到对应用户"
	UnableToQueryCommentErrorCode     = 50003
	UnableToQueryCommentError         = "无法查询到评论"
	UnableToCreateCommentErrorCode    = 50004
	UnableToCreateCommentError        = "无法创建评论"
	ActorIDNotMatchErrorCode          = 50005
	ActorIDNotMatchError              = "用户不匹配"
	UnableToDeleteCommentErrorCode    = 50006
	UnableToDeleteCommentError        = "无法删除评论"
	UnableToAddMessageErrorCode       = 50007
	UnableToAddMessageError           = "发送消息出错"
	UnableToQueryMessageErrorCode     = 50008
	UnableToQueryMessageError         = "查消息出错"
	PublishServiceInnerErrorCode      = 50009
	PublishServiceInnerError          = "发布服务出现内部错误，请稍后重试！"
	UnableToFollowErrorCode           = 50010
	UnableToFollowError               = "关注该用户失败"
	UnableToUnFollowErrorCode         = 50011
	UnableToUnFollowError             = "取消关注失败"
	UnableToGetFollowListErrorCode    = 50012
	UnableToGetFollowListError        = "无法查询到关注列表"
	UnableToGetFollowerListErrorCode  = 50013
	UnableToGetFollowerListError      = "无法查询到粉丝列表"
	UnableToRelateYourselfErrorCode   = 50014
	UnableToRelateYourselfError       = "无法关注自己"
	RelationNotFoundErrorCode         = 50015
	RelationNotFoundError             = "未关注该用户"
	StringToIntErrorCode              = 50016
	StringToIntError                  = "字符串转数字失败"
	RelationServiceIntErrorCode       = 50017
	RelationServiceIntError           = "关系服务出现内部错误"
	FavoriteServiceErrorCode          = 50018
	FavoriteServiceError              = "点赞服务内部出错"
	UserServiceInnerErrorCode         = 50019
	UserServiceInnerError             = "登录服务出现内部错误，请稍后重试！"
	UnableToQueryPollErrorCode        = 50020
	UnableToQueryPollError            = "无法查询到该提问"
	AlreadyFollowingErrorCode         = 50021
	AlreadyFollowingError             = "无法关注已关注的人"
	UnableToGetFriendListErrorCode    = 50022
	UnableToGetFriendListError        = "无法查询到好友列表"
	RecommendServiceInnerErrorCode    = 50023
	RecommendServiceInnerError        = "推荐系统内部错误"
	PollServiceInnerErrorCode         = 50024
	PollServiceInnerError             = "投票服务出现内部错误，请稍后重试！"
	UnableToCreateVoteErrorCode       = 50025
	UnableToCreateVoteError           = "无法创建投票"
	UnableToGetVoteCountListErrorCode = 50026
	UnableToGetVoteCountListError     = "无法查询到投票数"
	PollServiceFeedErrorCode          = 50027
	PollServiceFeedError              = "推送功能出错"
	NoUserUpdatedErrorCode            = 50028
	NoUserUpdatedError                = "没有用户被更新"
)

// Expected Error
const (
	AuthUserExistedCode           = 10001
	AuthUserExisted               = "用户名已存在！"
	UserNotExistedCode            = 10002
	UserNotExisted                = "用户不存在，请先注册或检查你的用户名是否正确！"
	AuthUserLoginFailedCode       = 10003
	AuthUserLoginFailed           = "用户信息错误，请检查账号密码是否正确"
	AuthUserNeededCode            = 10004
	AuthUserNeeded                = "用户无权限操作，请登陆后重试！"
	ActionCommentTypeInvalidCode  = 10005
	ActionCommentTypeInvalid      = "不合法的评论类型"
	ActionCommentLimitedCode      = 10006
	ActionCommentLimited          = "评论频繁，请稍后再试！"
	InvalidContentTypeCode        = 10007
	InvalidContentType            = "不合法的内容类型"
	FavoriteServiceDuplicateCode  = 10008
	FavoriteServiceDuplicateError = "不能重复点赞"
	FavoriteServiceCancelCode     = 10009
	FavoriteServiceCancelError    = "没有点赞,不能取消点赞"
	ChatActionLimitedCode         = 10010
	ChatActionLimitedError        = "发送消息频繁，请稍后再试！"
	FollowLimitedCode             = 10011
	FollowLimited                 = "关注频繁，请稍后再试！"
	UserDoNotExistedCode          = 10012
	UserDoNotExisted              = "查询用户不存在！"
)
