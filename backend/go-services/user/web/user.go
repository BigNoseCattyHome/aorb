package web

import (
	userModel "github.com/BigNoseCattyHome/aorb/backend/go-services/user/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/json"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var userClient user.UserServiceClient

func init() {
	userConn := grpc2.Connect(config.UserRpcServerName)
	userClient = user.NewUserServiceClient(userConn)
}

func UserHandler(c *gin.Context) {
	var req userModel.UserReq
	_, span := tracing.Tracer.Start(c.Request.Context(), "UserInfoHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("GateWay.UserInfo").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, userModel.UserRes{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
		})
		logging.SetSpanError(span, err)
		return
	}

	resp, err := userClient.GetUserInfo(c.Request.Context(), &user.UserRequest{
		UserId:  req.UserId,
		ActorId: req.ActorId,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when gateway get info from UserInfo Service")
		logging.SetSpanError(span, err)
		c.Render(http.StatusOK, json.CustomJSON{Data: resp, Context: c})
		return
	}

	c.Render(http.StatusOK, json.CustomJSON{Data: resp, Context: c})
}
