package report

import (
	"fmt"
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"time"
)

var (
	ReportLogic = newReportLogic()
)

type reportLogic struct{}

func newReportLogic() *reportLogic {
	return &reportLogic{}
}

type (
	ReqReport struct {
		UserId    int   `from:"userId" binding:"required"`
		StartTime int64 `from:"startTime binding" binding:"required"`
		EndTime   int64 `from:"endTime" binding:"required"`
	}

	RespReport struct {
		Projects []ProjectBasicReport `json:"projects"`
		Infos    []ProjectInfosReport `json:"infos"`
	}

	ProjectBasicReport struct {
		Id                  int     `json:"id"`
		Name                string  `json:"name"`
		Description         string  `json:"description"`
		Attr                int     `json:"attr"`
		Type                int     `json:"type"`
		State               int     `json:"state"`
		Capacity            float64 `json:"capacity"`
		Properties          string  `json:"properties"`
		Area                float64 `json:"area"`
		Star                int     `json:"star"`
		Address             string  `json:"address"`
		Connect             string  `json:"connect"`
		InvestmentAgreement string  `json:"investmentAgreement"`
		BusinessCondition   string  `json:"businessCondition"`
		BeginTime           string  `json:"beginTime"`
	}

	ProjectInfosReport struct {
		RecordTotal   int                  `json:"recordTotal"`
		AttachedTotal int                  `json:"attachedTotal"`
		Nodes         []ProjectInfosReport `json:"nodes"`
		Records       []ProjectInfosReport `json:"records"`
		Attached      []ProjectInfosReport `json:"attached"`
	}

	NodeReport struct {
		Id        int           `json:"id"`
		Node      int           `json:"node"`
		Name      string        `json:"name"`
		State     int           `json:"state"`
		CreatedAt string        `json:"createdAt"`
		UpdatedAT string        `json:"updatedAt"`
		Children  []*NodeReport `json:"children"`
	}

	RecordReport struct {
		Id          int    `json:"id"`
		ProjectId   int    `json:"projectId"`
		ProjectName string `json:"projectName"`
		NodeId      int    `json:"nodeId"`
		NodeName    string `json:"nodeName"`
		Overview    string `json:"overview"`
		State       int    `json:"state"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}

	AttachedReport struct {
		Id        int    `json:"id"`
		Url       string `json:"url"`
		Filename  string `json:"filename"`
		Mime      string `json:"mime"`
		Size      int64  `json:"size"`
		CreatedAt string `json:"createdAt"`
	}
)

func (rl *reportLogic) Report(req *ReqReport) (*RespReport, error) {
	//获取用户的所有项目
	projectIds, err := project.PersonLogic.GetUserProjectIds(req.UserId, false)
	if err != nil {
		return nil, err
	}

	projects, err := model.ProjectModel.SelectByIds(projectIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户项目出错")
	}
	fmt.Println(projects)

	startTime := time.Unix(req.StartTime, 0)
	endTime := time.Unix(req.EndTime, 0)

	//获取项目节点
	nodes, err := model.ProjectNodeModel.GetByProjectIds(projectIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目节点出错")
	}
	fmt.Println(nodes)

	//获取用户项目提交记录
	// userId , projectIds , createdAt range
	records, err := model.ProjectRecordModel.GetAllByProjectIdsAndUserId(projectIds, req.UserId, []*time.Time{
		&startTime, &endTime,
	})
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目记录出错")
	}
	fmt.Println(records)

	//获取用户项目上传文件记录
	attached, err := model.ProjectAttachedModel.GetAllByProjectIdsAndUserId(projectIds, req.UserId, []*time.Time{
		&startTime, &endTime,
	})
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目附件出错")
	}
	fmt.Println(attached)
	return nil, nil
}
