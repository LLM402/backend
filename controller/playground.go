package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
	"github.com/QuantumNous/new-api/middleware"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/types"

	"github.com/gin-gonic/gin"
)

func Playground(c *gin.Context) {
	var newAPIError *types.NewAPIError

	defer func() {
		if newAPIError != nil {
			c.JSON(newAPIError.StatusCode, gin.H{
				"error": newAPIError.ToOpenAIError(),
			})
		}
	}()

	useAccessToken := c.GetBool("use_access_token")
	if useAccessToken {
		newAPIError = types.NewError(errors.New("Access token is not currently supported."), types.ErrorCodeAccessDenied, types.ErrOptionWithSkipRetry())
		return
	}

	group := c.GetString("group")
	modelName := c.GetString("original_model")

	userId := c.GetInt("id")

	
	userCache, err := model.GetUserCache(userId)
	if err != nil {
		newAPIError = types.NewError(err, types.ErrorCodeQueryDataError, types.ErrOptionWithSkipRetry())
		return
	}
	userCache.WriteContext(c)

	tempToken := &model.Token{
		UserId: userId,
		Name:   fmt.Sprintf("playground-%s", group),
		Group:  group,
	}
	_ = middleware.SetupContextForToken(c, tempToken)
	_, newAPIError = getChannel(c, group, modelName, 0)
	if newAPIError != nil {
		return
	}
	
	common.SetContextKey(c, constant.ContextKeyRequestStartTime, time.Now())

	Relay(c, types.RelayFormatOpenAI)
}
