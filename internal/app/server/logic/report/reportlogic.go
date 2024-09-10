package report

import (
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic/project"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"fox_live_service/pkg/util/slicex"
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
		UserId    int     `form:"userId"`
		TimeRange []int64 `form:"timeRange[]" binding:"required"`
	}

	RespReport struct {
		Basic         *ProjectBasicReport `json:"basic"`
		RecordTotal   int                 `json:"recordTotal"`
		AttachedTotal int                 `json:"attachedTotal"`
		Nodes         []*NodeReport       `json:"nodes"`
		Records       []*RecordReport     `json:"records"`
		Attached      []*AttachedReport   `json:"attached"`
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

	NodeReport struct {
		Id       int           `json:"id"`
		NodeId   int           `json:"nodeId"`
		Name     string        `json:"name"`
		State    int           `json:"state"`
		Children []*NodeReport `json:"children,omitempty"`
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

func (rl *reportLogic) Report(req *ReqReport) ([]*RespReport, error) {
	//获取用户的所有项目
	projectIds, err := project.PersonLogic.GetUserProjectIds(req.UserId, false)
	if err != nil {
		return nil, err
	}

	startTime := time.UnixMilli(req.TimeRange[0])
	endTime := time.UnixMilli(req.TimeRange[1])

	//获取用户这个时间段内又提交记录的project
	projectRecords, err := model.ProjectRecordModel.GetRecordByUserIdAndTimeRange(req.UserId, startTime, endTime)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户提交记录出错")
	}

	projectRecordIds := make([]int, 0, len(projectRecords))
	for _, projectRecord := range projectRecords {
		projectRecordIds = append(projectRecordIds, projectRecord.ProjectId)
	}

	projectIds = slicex.IntersectionInt(projectIds, projectRecordIds)

	projects, err := model.ProjectModel.SelectByIds(projectIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户项目出错")
	}

	basicProjects := make([]*ProjectBasicReport, 0, len(projects))
	for _, p := range projects {
		basicProjects = append(basicProjects, &ProjectBasicReport{
			Id:                  p.Id,
			Name:                p.Name,
			Description:         p.Description,
			Attr:                p.Attr,
			Type:                p.Type,
			State:               p.State,
			Capacity:            p.Capacity,
			Properties:          p.Properties,
			Area:                p.Area,
			Star:                p.Star,
			Address:             p.Address,
			Connect:             p.Connect,
			InvestmentAgreement: p.InvestmentAgreement,
			BusinessCondition:   p.BusinessCondition,
			BeginTime:           p.BeginTime.Format(global.TimeFormat),
		})
	}

	//获取项目节点
	nodes, err := model.ProjectNodeModel.GetByProjectIds(projectIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目节点出错")
	}
	fmt.Println(len(nodes))
	nodeMap := make(map[int][]*model.ProjectNode)
	for _, v := range nodes {
		nodeMap[v.ProjectId] = append(nodeMap[v.ProjectId], v)
	}

	//获取用户项目提交记录
	records, err := model.ProjectRecordModel.GetAllByProjectIdsAndUserId(projectIds, req.UserId, []*time.Time{
		&startTime, &endTime,
	})
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目记录出错")
	}
	recordMap := make(map[int][]*RecordReport)
	for _, v := range records {
		recordMap[v.ProjectId] = append(recordMap[v.ProjectId], &RecordReport{
			Id:        v.Id,
			NodeName:  v.NodeName,
			Overview:  v.Overview,
			State:     v.State,
			NodeId:    v.NodeId,
			CreatedAt: v.CreatedAt.Format(global.TimeFormat),
		})
	}

	//获取用户项目上传文件记录
	attached, err := model.ProjectAttachedModel.GetAllByProjectIdsAndUserId(projectIds, req.UserId, []*time.Time{
		&startTime, &endTime,
	})
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目附件出错")
	}
	attachedMap := make(map[int][]*AttachedReport)
	for _, v := range attached {
		attachedMap[v.ProjectId] = append(attachedMap[v.ProjectId], &AttachedReport{
			Id:        v.Id,
			Url:       v.Url,
			Filename:  v.Filename,
			Mime:      v.Mime,
			Size:      v.Size,
			CreatedAt: v.CreatedAt.Format(global.TimeFormat),
		})
	}

	res := make([]*RespReport, 0, len(projectIds))

	for _, v := range basicProjects {
		res = append(res, rl.genReportItem(v, nodeMap, recordMap, attachedMap))
	}

	return res, nil
}

func (rl *reportLogic) genReportItem(basis *ProjectBasicReport, nodeMap map[int][]*model.ProjectNode,
	recordMap map[int][]*RecordReport, attachedMap map[int][]*AttachedReport) *RespReport {

	records := make([]*RecordReport, 0)
	if _, ok := recordMap[basis.Id]; ok {
		records = recordMap[basis.Id]
	}

	nodeStateMap := make(map[int]int)
	for _, v := range records {
		if _, ok := nodeStateMap[v.NodeId]; !ok {
			if v.State == model.ProjectRecordStateFinished {
				nodeStateMap[v.NodeId] = 99
			}
		}
	}

	attached := make([]*AttachedReport, 0)
	if _, ok := attachedMap[basis.Id]; ok {
		attached = attachedMap[basis.Id]
	}

	return &RespReport{
		Basic:         basis,
		RecordTotal:   len(records),
		AttachedTotal: len(attached),
		Nodes:         rl.formatNodes(nodeMap[basis.Id], nodeStateMap),
		Records:       records,
		Attached:      attached,
	}
}

func (rl *reportLogic) formatNodes(nodes []*model.ProjectNode, nodeState map[int]int) []*NodeReport {
	pNodeMap := make(map[int][]*NodeReport)
	for _, v := range nodes {
		if _, ok := pNodeMap[v.PId]; !ok {
			pNodeMap[v.PId] = make([]*NodeReport, 0)
		}
		state := v.State
		if _, ok := nodeState[v.PId]; ok {
			state = nodeState[v.PId]
		}
		pNodeMap[v.PId] = append(pNodeMap[v.PId], &NodeReport{
			Id:     v.Id,
			NodeId: v.NodeId,
			Name:   v.Name,
			State:  state,
		})
	}

	res := make([]*NodeReport, 0, len(pNodeMap[0]))
	for _, node := range pNodeMap[0] {
		node.Children = pNodeMap[node.NodeId]
		res = append(res, node)
	}
	fmt.Println(len(res))

	return res
}
