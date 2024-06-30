package middleware

import (
	"github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
	"strconv"
)

var client auth.AuthServiceClient

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := tracing.Tracer.Start(c.Request.Context(), "TokenAuthMiddleware")
		defer span.End()
		logging.SetSpanWithHostname(span)
		logger := logging.LogService("GateWay.AuthMiddleWare").WithContext(ctx)
		span.SetAttributes(attribute.String("url", c.Request.URL.Path))
		if c.Request.URL.Path == "/aorb/user/login/" ||
			c.Request.URL.Path == "/aorb/user/register/" ||
			c.Request.URL.Path == "/aorb/comment/list/" ||
			c.Request.URL.Path == "/aorb/poll/list/" ||
			c.Request.URL.Path == "/ping" {
			// 直接放行，这里后期需要修改
			c.Request.URL.RawQuery += "&actor_id=" + config.Conf.Other.AnonymityUser
			span.SetAttributes(attribute.String("make_url", c.Request.URL.String()))
			logger.WithFields(logrus.Fields{
				"Path": c.Request.URL.Path,
			}).Debugf("Skip Auth with targeted url")
			c.Next()
			return
		}

		//var token string
		//if c.Request.URL.Path == "/douyin/publish/action/" {
		//	token = c.PostForm("token")
		//} else {
		//	token = c.Query("token")
		//}

		token := c.Query("token")
		if token == "" && (c.Request.URL.Path == "/aorb/feed/" ||
			c.Request.URL.Path == "/aorb/relation/follow/list/" ||
			c.Request.URL.Path == "/aorb/relation/follower/list/") {
			c.Request.URL.RawQuery += "&actor_id=" + config.Conf.Other.AnonymityUser
			span.SetAttributes(attribute.String("mark_url", c.Request.URL.String()))
			logger.WithFields(logrus.Fields{
				"Path": c.Request.URL.Path,
			}).Debugf("Skip Auth with targeted url")
			c.Next()
			return
		}
		span.SetAttributes(attribute.String("token", token))

		// 核实用户token
		authenticate, err := client.Authenticate(c.Request.Context(), &auth.AuthenticateRequest{
			Token: token,
		})
		if err != nil {
			logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Gateway Auth meet trouble")
			span.RecordError(err)
			c.JSON(http.StatusOK, gin.H{
				"status_code": strings.GateWayErrorCode,
				"status_msg":  strings.GateWayError,
			})
			c.Abort()
			return
		}

		if authenticate.StatusCode != 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status_code": strings.AuthUserNeededCode,
				"status_msg":  strings.AuthUserNeeded,
			})
			c.Abort()
			return
		}

		c.Request.URL.RawQuery += "&actor_id=" + strconv.FormatUint(uint64(authenticate.UserId), 10)
		c.Next()
	}
}

func init() {
	authConn := grpc2.Connect(config.AuthRpcServerName)
	client = auth.NewAuthServiceClient(authConn)
}
