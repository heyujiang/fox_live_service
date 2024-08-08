package handler

import (
	"fox_live_service/internal/app/server/logic/system"
	"fox_live_service/pkg/common"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
)

var RuleHandler = newRuleHandler()

type ruleHandler struct{}

func newRuleHandler() *ruleHandler {
	return &ruleHandler{}
}

func (h *ruleHandler) Create(c *gin.Context) {
	var req system.ReqCreateRule
	if err := c.ShouldBindJSON(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := system.RuleLogic.Create(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *ruleHandler) Delete(c *gin.Context) {
	var req system.ReqDeleteRule
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := system.RuleLogic.Delete(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *ruleHandler) Update(c *gin.Context) {
	var reqUri system.ReqUpdateRuleUri
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody system.ReqUpdateRuleBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := system.ReqUpdateRule{
		ReqUpdateRuleUri:  reqUri,
		ReqUpdateRuleBody: reqBody,
	}

	res, err := system.RuleLogic.Update(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *ruleHandler) List(c *gin.Context) {
	res, err := system.RuleLogic.List()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *ruleHandler) Parents(c *gin.Context) {
	res, err := system.RuleLogic.Parents()
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *ruleHandler) UpdateStatus(c *gin.Context) {
	var reqUri system.ReqUpdateRuleUri
	if err := c.ShouldBindUri(&reqUri); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}
	var reqBody system.ReqUpdateRuleStatusBody
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	req := system.ReqUpdateRuleStatus{
		ReqUpdateRuleUri:        reqUri,
		ReqUpdateRuleStatusBody: reqBody,
	}

	res, err := system.RuleLogic.UpdateStatus(&req, c.GetInt("uid"))
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}

func (h *ruleHandler) GetRoleRules(c *gin.Context) {
	var req system.ReqGetRoleRules
	if err := c.ShouldBindUri(&req); err != nil {
		common.ResponseErr(c, errorx.NewErrorX(errorx.ErrParam, "param error"))
		return
	}

	res, err := system.RuleLogic.GetRules(&req)
	if err != nil {
		common.ResponseErr(c, err)
		return
	}

	common.ResponseOK(c, res)
	return
}
