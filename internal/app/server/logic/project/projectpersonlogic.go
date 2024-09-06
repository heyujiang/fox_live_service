package project

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
)

var PersonLogic = newPersonLogic()

type (
	personLogic struct{}

	ReqCreateProjectPerson struct {
		ProjectId int `json:"projectId"`
		UserId    int `json:"userId"`
		Type      int `json:"type"`
	}

	RespCreateProjectPerson struct{}

	ReqDeleteProjectPerson struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectPerson struct{}

	ReqUpdateProjectPerson struct {
		ReqUriUpdateProjectPerson
		ReqBodyUpdateProjectPerson
	}

	ReqUriUpdateProjectPerson struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectPerson struct {
	}

	RespUpdateProjectPerson struct {
	}

	ReqInfoProjectPerson struct {
		Id int `uri:"id"`
	}

	RespInfoProjectPerson struct {
	}

	ReqProjectPersonList struct {
		ProjectId int `uri:"projectId"`
	}

	RespProjectPersonList struct {
		List []*ListProjectPersonItem
	}

	ListProjectPersonItem struct {
		Id          int    `json:"id"`
		ProjectId   int    `json:"project_id"`
		UserId      int    `json:"user_id"`
		Type        int    `json:"type"`
		Name        string `json:"name"`
		PhoneNumber string `json:"phone_number"`
		CreatedId   int    `json:"created_id"`
		UpdatedId   int    `json:"updated_id"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
)

func newPersonLogic() *personLogic {
	return &personLogic{}
}

func (b *personLogic) Create(req *ReqCreateProjectPerson, uid int) (*RespCreateProjectPerson, error) {

	if req.Type == model.ProjectPersonTypeFirst {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "不能添加第一负责人")
	}

	user, err := model.UserModel.Find(req.UserId)
	if err != nil {
		slog.Error("create project person get user error ", "id", req.UserId, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "用户不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户错误")
	}

	project, err := model.ProjectModel.Find(req.ProjectId)
	if err != nil {
		slog.Error("create project person get project error ", "id", req.UserId, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目错误")
	}

	//查询用户是否为项目成员
	projectPerson, err := model.ProjectPersonModel.FindByProjectIdAndUserId(req.ProjectId, req.UserId)
	if err != nil && !errors.Is(err, model.ErrNotRecord) {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目成员失败")
	}
	if projectPerson != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "用户已经是项目成员了")
	}

	if err := model.ProjectPersonModel.Create(&model.ProjectPerson{
		ProjectId:   project.Id,
		UserId:      user.Id,
		Name:        user.Username,
		PhoneNumber: user.PhoneNumber,
		Type:        req.Type,
		CreatedId:   uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目成员失败")
	}

	return &RespCreateProjectPerson{}, nil
}

func (b *personLogic) Delete(req *ReqDeleteProjectPerson, uid int) (*RespDeleteProjectPerson, error) {
	if err := model.ProjectPersonModel.Delete(req.Id, uid); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除项目成员失败")
	}

	return &RespDeleteProjectPerson{}, nil
}

func (b *personLogic) List(req *ReqProjectPersonList) ([]*ListProjectPersonItem, error) {
	persons, err := model.ProjectPersonModel.SelectByProjectId(req.ProjectId)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目成员失败")
	}
	res := make([]*ListProjectPersonItem, 0, len(persons))
	for _, v := range persons {
		res = append(res, &ListProjectPersonItem{
			Id:          v.Id,
			ProjectId:   v.ProjectId,
			UserId:      v.UserId,
			Type:        v.Type,
			Name:        v.Name,
			PhoneNumber: v.PhoneNumber,
			CreatedId:   v.CreatedId,
			CreatedAt:   v.CreatedAt.Format(global.TimeFormat),
		})
	}

	return res, nil
}

// GetUserProjectIds 获取用户所有项目的id
// isFirst 是否只获取作为第一负责人的项目
func (b *personLogic) GetUserProjectIds(userId int, isFirst bool) ([]int, error) {
	projectPersons := make([]*model.ProjectPerson, 0)
	var err error
	if isFirst {
		projectPersons, err = model.ProjectPersonModel.SelectByUserIdAndFirst(userId)
	} else {
		projectPersons, err = model.ProjectPersonModel.SelectByUserId(userId)
	}
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户项目id失败")
	}

	projectIds := make([]int, 0, len(projectPersons))
	for _, v := range projectPersons {
		projectIds = append(projectIds, v.ProjectId)
	}
	return projectIds, nil
}
