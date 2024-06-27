package web

import (
	authModels "github.com/BigNoseCattyHome/aorb/backend/go-services/auth/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
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

var Client auth.AuthServiceClient

func LoginHandler(c *gin.Context) {
	var req authModels.LoginRequest
	_, span := tracing.Tracer.Start(c.Request.Context(), "LoginHandler")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("GateWay.Login").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, authModels.LoginResponse{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
			UserId:     "0",
			Token:      "",
		})
		return
	}

	res, err := Client.Login(c.Request.Context(), &auth.LoginRequest{
		Username: req.UserName,
		Password: req.Password,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Warnf("Error when trying to connect with AuthService")
		c.Render(http.StatusOK, json.CustomJSON{Data: res, Context: c})
		return
	}

	logger.WithFields(logrus.Fields{
		"Username": req.UserName,
		"Token":    res.Token,
		"UserId":   res.UserId,
	}).Infof("User log in")

	c.Render(http.StatusOK, json.CustomJSON{Data: res, Context: c})
}

func RegisterHandler(c *gin.Context) {
	var req authModels.RegisterRequest
	_, span := tracing.Tracer.Start(c.Request.Context(), "LoginHandler")
	defer span.End()
	logger := logging.LogService("GateWay.Register").WithContext(c.Request.Context())

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, authModels.RegisterResponse{
			StatusCode: strings.GateWayParamsErrorCode,
			StatusMsg:  strings.GateWayParamsError,
			UserId:     "0",
			Token:      "",
		})
		return
	}

	res, err := Client.Register(c.Request.Context(), &auth.RegisterRequest{
		Username: req.UserName,
		Password: req.Password,
	})

	if err != nil {
		logger.WithFields(logrus.Fields{
			"Username": req.UserName,
		}).Warnf("Error when trying to connect with AuthService")
		c.Render(http.StatusOK, json.CustomJSON{Data: res, Context: c})
		return
	}

	logger.WithFields(logrus.Fields{
		"Username": req.UserName,
		"Token":    res.Token,
		"UserId":   res.UserId,
	}).Infof("User register in")

	c.Render(http.StatusOK, json.CustomJSON{Data: res, Context: c})
}

func init() {
	conn := grpc2.Connect(config.AuthRpcServerName)
	Client = auth.NewAuthServiceClient(conn)
}
