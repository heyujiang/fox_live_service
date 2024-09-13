package project

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
)

var AttachedLogic = newAttachedLogic()

type (
	attachedLogic struct{}

	ReqCreateProjectAttached struct {
		RecordId int    `json:"recordId"`
		Filename string `json:"filename"`
		Url      string `json:"url"`
		Mime     string `json:"mime"`
		Size     int64  `json:"size"`
	}

	RespCreateProjectAttached struct{}

	ReqDeleteProjectAttached struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectAttached struct{}

	ReqUpdateProjectAttached struct {
		ReqUriUpdateProjectAttached
		ReqBodyUpdateProjectAttached
	}

	ReqUriUpdateProjectAttached struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectAttached struct {
	}

	RespUpdateProjectAttached struct {
	}

	ReqInfoProjectAttached struct {
		Id int `uri:"id"`
	}

	RespInfoProjectAttached struct {
	}

	ReqProjectAttachedList struct {
		ProjectId int `form:"projectId"`
		NodeId    int `form:"nodeId"`
		RecordId  int `form:"recordId"`
		UserId    int `form:"userId"`
	}

	RespProjectAttachedList struct {
		Id        int    `json:"id"`
		ProjectId int    `json:"projectId"`
		NodeId    int    `json:"nodeId"`
		RecordId  int    `json:"recordId"`
		UserId    int    `json:"userId"`
		Url       string `json:"url"`
		Filename  string `json:"filename"`
		Mime      string `json:"mime"`
		Size      int64  `json:"size"`
		CreatedAt string `json:"createdAt"`
	}

	AttachedFile struct {
		Url      string `json:"url"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
		Mime     string `json:"mime"`
	}
)

func newAttachedLogic() *attachedLogic {
	return &attachedLogic{}
}

func (b *attachedLogic) Create(req *ReqCreateProjectAttached, uid int) (*RespCreateProjectAttached, error) {
	record, err := model.ProjectRecordModel.Find(req.RecordId)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "记录不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询记录出错")
	}

	_, err = model.ProjectModel.Find(record.ProjectId)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目出错")
	}

	hasProject, err := PersonLogic.CheckUserHasProject(uid, record.ProjectId)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, err.Error())
	}
	if !hasProject {
		slog.Error("不属于当前项目的项目成员.", "projectId", record.ProjectId, "userId", uid)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "不是项目成员，不能上传项目附件")
	}

	if err := model.ProjectAttachedModel.Create(&model.ProjectAttached{
		ProjectId: record.ProjectId,
		NodeId:    record.NodeId,
		RecordId:  record.Id,
		UserId:    uid,
		Url:       req.Url,
		Filename:  req.Filename,
		Mime:      req.Mime,
		Size:      req.Size,
		CreatedId: uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "保存附件出错")
	}

	return &RespCreateProjectAttached{}, nil
}

func (b *attachedLogic) Delete(req *ReqDeleteProjectAttached) (*RespDeleteProjectAttached, error) {
	return &RespDeleteProjectAttached{}, nil
}

func (b *attachedLogic) Update(req *ReqUpdateProjectAttached) (*RespUpdateProjectAttached, error) {
	return &RespUpdateProjectAttached{}, nil
}

func (b *attachedLogic) Info(req *ReqInfoProjectAttached) (*RespInfoProjectAttached, error) {
	return &RespInfoProjectAttached{}, nil
}

func (b *attachedLogic) List(req *ReqProjectAttachedList) ([]*RespProjectAttachedList, error) {
	attenheds, err := model.ProjectAttachedModel.GetProjectAttachedByCond(b.getProjectAttachedCond(req))
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询附件出错")
	}

	res := make([]*RespProjectAttachedList, 0, len(attenheds))
	for _, v := range attenheds {
		res = append(res, &RespProjectAttachedList{
			Id:        v.Id,
			ProjectId: v.ProjectId,
			NodeId:    v.NodeId,
			RecordId:  v.RecordId,
			UserId:    v.UserId,
			Url:       v.Url,
			Filename:  v.Filename,
			Mime:      v.Mime,
			Size:      v.Size,
			CreatedAt: v.CreatedAt.Format(global.TimeFormat),
		})
	}
	return res, nil
}

func (b *attachedLogic) getProjectAttachedCond(req *ReqProjectAttachedList) *model.ProjectAttachedCond {
	cond := &model.ProjectAttachedCond{}
	if req == nil {
		return cond
	}
	if req.ProjectId > 0 {
		cond.ProjectId = req.ProjectId
	}

	if req.NodeId > 0 {
		cond.NodeId = req.NodeId
	}

	if req.RecordId > 0 {
		cond.RecordId = req.RecordId
	}

	if req.UserId > 0 {
		cond.UserId = req.UserId
	}

	return cond
}
