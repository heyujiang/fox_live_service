package project

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
)

var ContactLogic = newContactLogic()

type (
	contactLogic struct{}

	ReqCreateProjectContact struct {
		ProjectId   int    `json:"projectId"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		Type        int    `json:"type"`
		Description string `json:"description"`
	}

	RespCreateProjectContact struct{}

	ReqDeleteProjectContact struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectContact struct{}

	ReqUpdateProjectContact struct {
		ReqUriUpdateProjectContact
		ReqBodyUpdateProjectContact
	}

	ReqUriUpdateProjectContact struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectContact struct {
	}

	RespUpdateProjectContact struct {
	}

	ReqInfoProjectContact struct {
		Id int `uri:"id"`
	}

	RespInfoProjectContact struct {
	}

	ReqProjectContactList struct {
		ProjectId int `uri:"projectId"`
	}

	RespProjectContactList struct {
		List []*ListProjectContactItem
	}

	ListProjectContactItem struct {
		Id          int    `json:"id"`
		ProjectId   int    `json:"projectId"`
		Type        int    `json:"type"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		Description string `json:"description"`
		CreatedId   int    `json:"createdId"`
		CreatedAt   string `json:"createdAt"`
	}
)

func newContactLogic() *contactLogic {
	return &contactLogic{}
}

func (b *contactLogic) Create(req *ReqCreateProjectContact, uid int) (*RespCreateProjectContact, error) {

	project, err := model.ProjectModel.Find(req.ProjectId)
	if err != nil {
		slog.Error("create project contact get project error ", "id", req.ProjectId, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目错误")
	}

	if err := model.ProjectContactModel.Create(&model.ProjectContact{
		ProjectId:   project.Id,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Type:        req.Type,
		Description: req.Description,
		CreatedId:   uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目联系失败")
	}

	return &RespCreateProjectContact{}, nil
}

func (b *contactLogic) Delete(req *ReqDeleteProjectContact) (*RespDeleteProjectContact, error) {
	if err := model.ProjectContactModel.Delete(req.Id); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除项目联系人失败")
	}

	return &RespDeleteProjectContact{}, nil
}

func (b *contactLogic) List(req *ReqProjectContactList) (*RespProjectContactList, error) {
	Contacts, err := model.ProjectContactModel.SelectByProjectId(req.ProjectId)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目联系人失败")
	}
	res := make([]*ListProjectContactItem, 0, len(Contacts))
	for _, v := range Contacts {
		res = append(res, &ListProjectContactItem{
			Id:          v.Id,
			ProjectId:   v.ProjectId,
			Type:        v.Type,
			Name:        v.Name,
			PhoneNumber: v.PhoneNumber,
			Description: v.Description,
			CreatedId:   v.CreatedId,
			CreatedAt:   v.CreatedAt.Format(global.TimeFormat),
		})
	}

	return &RespProjectContactList{
		List: res,
	}, nil
}
