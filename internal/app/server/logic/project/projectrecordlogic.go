package project

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
	"time"
)

var RecordLogic = newRecordLogic()

type (
	recordLogic struct{}

	ReqCreateProjectRecord struct {
		ProjectId int           `json:"projectId"`
		NodeId    int           `json:"nodeId"`
		Overview  string        `json:"overview"`
		State     int           `json:"state"`
		File      *AttachedFile `json:"file"`
	}

	RespCreateProjectRecord struct{}

	ReqDeleteProjectRecord struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectRecord struct{}

	ReqUpdateProjectRecord struct {
		ReqUriUpdateProjectRecord
		ReqBodyUpdateProjectRecord
	}

	ReqUriUpdateProjectRecord struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectRecord struct {
		Overview string `json:"overview"`
	}

	RespUpdateProjectRecord struct {
	}

	ReqInfoProjectRecord struct {
		Id int `uri:"id"`
	}

	RespInfoProjectRecord struct {
	}

	ReqProjectRecordList struct {
		logic.ReqPage
		ReqFromProjectRecordList
	}

	ReqFromProjectRecordList struct {
		ProjectId int `form:"projectId"`
		NodeId    int `form:"nodeId"`
		UserId    int `form:"userId"`
	}

	RespProjectRecordList struct {
		Page  int                      `json:"page"`
		Size  int                      `json:"size"`
		Count int                      `json:"count"`
		List  []*ListProjectRecordItem `json:"list"`
	}

	ListProjectRecordItem struct {
		Id          int    `json:"id"`
		ProjectId   int    `json:"projectId"`
		ProjectName string `json:"projectName"`
		NodeId      int    `json:"nodeId"`
		NodeName    string `json:"nodeName"`
		UserId      int    `json:"userId"`
		Username    string `json:"username"`
		Overview    string `json:"overview"`
		State       int    `json:"state"`
		CreatedAt   string `json:"createdAt"`
		UpdatedAt   string `json:"updatedAt"`
	}
)

func newRecordLogic() *recordLogic {
	return &recordLogic{}
}

func (b *recordLogic) Create(req *ReqCreateProjectRecord, uid int, username string) (*RespCreateProjectRecord, error) {
	//查询项目
	project, err := model.ProjectModel.Find(req.ProjectId)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目错误")
	}

	projectNode, err := model.ProjectNodeModel.FindByProjectIdAndNodeId(req.ProjectId, req.NodeId)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目节点不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目节点错误")
	}

	projectRecord := &model.ProjectRecord{
		ProjectId:   req.ProjectId,
		ProjectName: project.Name,
		NodeId:      req.NodeId,
		NodeName:    projectNode.Name,
		UserId:      uid,
		Username:    username,
		Overview:    req.Overview,
		State:       req.State,
		CreatedId:   uid,
		UpdatedId:   uid,
	}
	if err := model.ProjectRecordModel.Create(projectRecord); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "新增项目记录出错")
	}

	if req.State == model.ProjectRecordStateFinished { // 修改项目节点为完成
		if err := model.ProjectNodeModel.UpdateProjectNodeState(projectNode.Id, model.ProjectNodeStateFinished, uid); err != nil {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "修改项目状态出错")
		}
	} else if req.State == model.ProjectRecordStateIng { //进行中
		if err := model.ProjectNodeModel.UpdateProjectNodeState(projectNode.Id, model.ProjectNodeStateInProcess, uid); err != nil {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "修改项目状态出错")
		}
	}

	if req.File != nil {
		attached := model.ProjectAttached{
			ProjectId: req.ProjectId,
			NodeId:    req.NodeId,
			RecordId:  projectRecord.Id,
			UserId:    uid,
			Url:       req.File.Url,
			Filename:  req.File.Filename,
			Mime:      req.File.Mime,
			Size:      req.File.Size,
			CreatedId: 0,
			CreatedAt: time.Time{},
		}

		if err := model.ProjectAttachedModel.Create(&attached); err != nil {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "保存附件信息出错")
		}
	}

	return &RespCreateProjectRecord{}, nil
}

func (b *recordLogic) Delete(req *ReqDeleteProjectRecord) (*RespDeleteProjectRecord, error) {
	return &RespDeleteProjectRecord{}, nil
}

func (b *recordLogic) Update(req *ReqUpdateProjectRecord, uid int) (*RespUpdateProjectRecord, error) {
	//查询项目
	project, err := model.ProjectRecordModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目记录不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目记录错误")
	}

	if err := model.ProjectRecordModel.Update(&model.ProjectRecord{
		Id:        project.Id,
		Overview:  req.Overview,
		UpdatedId: uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "编辑项目记录出错")
	}

	return &RespUpdateProjectRecord{}, nil
}

func (b *recordLogic) Info(req *ReqInfoProjectRecord) (*RespInfoProjectRecord, error) {
	return &RespInfoProjectRecord{}, nil
}

func (b *recordLogic) List(req *ReqProjectRecordList) (*RespProjectRecordList, error) {
	logic.VerifyReqPage(&req.ReqPage)
	cond := b.buildSearchCond(req)
	totalCount, err := model.ProjectRecordModel.GetProjectRecordCountByCond(cond)
	if err != nil {
		slog.Error("list project record count error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目记录列表错误")
	}
	projects, err := model.ProjectRecordModel.GetProjectRecordByCond(cond, req.Page, req.Size)
	if err != nil {
		slog.Error("list project record list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目记录列表错误")
	}

	items := make([]*ListProjectRecordItem, 0, len(projects))
	for _, pro := range projects {
		items = append(items, &ListProjectRecordItem{
			Id:          pro.Id,
			ProjectId:   pro.ProjectId,
			ProjectName: pro.ProjectName,
			NodeId:      pro.NodeId,
			NodeName:    pro.NodeName,
			UserId:      pro.UserId,
			Username:    pro.Username,
			Overview:    pro.Overview,
			State:       pro.State,
			CreatedAt:   pro.CreatedAt.Format(global.TimeFormat),
			UpdatedAt:   pro.UpdatedAt.Format(global.TimeFormat),
		})
	}

	return &RespProjectRecordList{
		Page:  req.Page,
		Size:  req.Size,
		Count: totalCount,
		List:  items,
	}, nil
}

func (b *recordLogic) ListNoPage(req *ReqProjectRecordList) ([]*ListProjectRecordItem, error) {
	cond := b.buildSearchCond(req)
	req.Page = 1
	req.Size = 1000
	projects, err := model.ProjectRecordModel.GetProjectRecordByCond(cond, req.Page, req.Size)
	if err != nil {
		slog.Error("list project record list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目记录列表错误")
	}

	items := make([]*ListProjectRecordItem, 0, len(projects))
	for _, pro := range projects {
		items = append(items, &ListProjectRecordItem{
			Id:          pro.Id,
			ProjectId:   pro.ProjectId,
			ProjectName: pro.ProjectName,
			NodeId:      pro.NodeId,
			NodeName:    pro.NodeName,
			UserId:      pro.UserId,
			Username:    pro.Username,
			Overview:    pro.Overview,
			State:       pro.State,
			CreatedAt:   pro.CreatedAt.Format(global.TimeFormat),
			UpdatedAt:   pro.UpdatedAt.Format(global.TimeFormat),
		})
	}

	return items, nil
}

func (b *recordLogic) buildSearchCond(req *ReqProjectRecordList) *model.ProjectRecordCond {
	cond := &model.ProjectRecordCond{}

	if req.ProjectId > 0 {
		cond.ProjectId = req.ProjectId
	}

	if req.UserId != 0 {
		cond.UserId = req.UserId
	}

	if req.NodeId != 0 {
		cond.NodeId = req.NodeId
	}

	return cond
}

//func ()  {

//}
