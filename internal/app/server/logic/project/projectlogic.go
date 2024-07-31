package project

import (
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic"
	"fox_live_service/internal/app/server/logic/node"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
	"time"
)

var BisLogic = newBisLogic()

type (
	bisLogic struct{}

	ReqCreateProject struct {
		Name                string                  `json:"name"`
		Description         string                  `json:"description"`
		Attr                int                     `json:"attr"`
		Type                int                     `json:"type"`
		State               int                     `json:"state"`
		Capacity            float64                 `json:"capacity"`
		Properties          string                  `json:"properties"`
		Area                float64                 `json:"area"`
		Star                int                     `json:"star"`
		Address             string                  `json:"address"`
		Connect             string                  `json:"connect"`
		InvestmentAgreement string                  `json:"investmentAgreement"`
		BusinessCondition   string                  `json:"businessCondition"`
		BeginTime           int64                   `json:"beginTime"`
		Contact             []*CreateProjectContact `json:"contact"`
		Person              []*CreateProjectPerson  `json:"person"`
		NodeIds             []int                   `json:"nodeIds"`
	}

	CreateProjectContact struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		Type        int    `json:"type"`
		Description string `json:"description"`
	}

	CreateProjectPerson struct {
		UserId int `json:"userId"`
		Type   int `json:"type"`
	}

	RespCreateProject struct{}

	ReqDeleteProject struct {
		Id int `uri:"id"`
	}

	RespDeleteProject struct{}

	ReqUpdateProject struct {
		ReqUriUpdateProject
		ReqBodyUpdateProject
	}

	ReqUriUpdateProject struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProject struct {
		Name                string  `json:"name"`
		Description         string  `json:"description"`
		Attr                int     `json:"attr"`
		State               int     `json:"state"`
		Type                int     `json:"type"`
		Capacity            float64 `json:"capacity"`
		Properties          string  `json:"properties"`
		Area                float64 `json:"area"`
		Address             string  `json:"address"`
		Connect             string  `json:"connect"`
		InvestmentAgreement string  `json:"investmentAgreement"`
		BusinessCondition   string  `json:"businessCondition"`
		BeginTime           int64   `json:"beginTime"`
	}

	RespUpdateProject struct {
	}

	ReqInfoProject struct {
		Id int `uri:"id"`
	}

	RespInfoProject struct {
		Id                  int       `json:"id"`
		Name                string    `json:"name"`
		Description         string    `json:"description"`
		Attr                int       `json:"attr"`
		State               int       `json:"state"`
		Type                int       `json:"type"`
		NodeId              int       `json:"nodeId"`
		NodeName            string    `json:"nodName"`
		Schedule            float64   `json:"schedule"`
		Capacity            float64   `json:"capacity"`
		Properties          string    `json:"properties"`
		Area                float64   `json:"area"`
		Address             string    `json:"address"`
		Connect             string    `json:"connect"`
		InvestmentAgreement string    `json:"investmentAgreement"`
		BusinessCondition   string    `json:"businessCondition"`
		BeginTime           time.Time `json:"beginTime"`
		CreatedId           int       `json:"createdId"`
		UpdatedId           int       `json:"updatedId"`
		CreatedAt           string    `json:"createdAt"`
		UpdatedAt           string    `json:"updatedAt"`
	}

	ReqProjectList struct {
		logic.ReqPage
		ReqFromProjectList
	}

	ReqFromProjectList struct {
		Id     int    `form:"id"`
		Name   string `form:"name"`
		NodeId int    `form:"node_id"`
		Attr   int    `form:"attr"`
		State  int    `form:"state"`
		Type   int    `form:"type"`
	}

	RespProjectList struct {
		Page  int                `json:"page"`
		Size  int                `json:"size"`
		Count int                `json:"count"`
		List  []*ListProjectItem `json:"list"`
	}

	ListProjectItem struct {
		Id                  int     `json:"id"`
		Name                string  `json:"name"`
		Attr                int     `json:"attr"`
		State               int     `json:"state"`
		Type                int     `json:"type"`
		NodeName            string  `json:"node_name"`
		Schedule            float64 `json:"schedule"`
		Capacity            float64 `json:"capacity"`
		Properties          string  `json:"properties"`
		Area                float64 `json:"area"`
		Address             string  `json:"address"`
		Connect             string  `json:"connect"`
		InvestmentAgreement string  `json:"investmentAgreement"`
		BusinessCondition   string  `json:"businessCondition"`
		BeginTime           string  `json:"beginTime"`
		CreatedAt           string  `json:"createdAt"`
		UpdatedAt           string  `json:"updatedAt"`
	}
)

func newBisLogic() *bisLogic {
	return &bisLogic{}
}

func (b *bisLogic) Create(req *ReqCreateProject, uid int) (*RespCreateProject, error) {
	projectNodes, err := b.buildProjectNode(req.NodeIds, uid)
	projectContacts, err := b.buildProjectContact(req.Contact, uid)
	projectPersons, err := b.buildProjectPerson(req.Person, uid)

	projectId, err := model.ProjectModel.Create(&model.Project{
		Name:                req.Name,
		Description:         req.Description,
		Attr:                req.Attr,
		State:               req.State,
		Type:                req.Type,
		Capacity:            req.Capacity,
		Properties:          req.Properties,
		Area:                req.Area,
		Address:             req.Address,
		Connect:             req.Connect,
		Star:                req.Star,
		InvestmentAgreement: req.InvestmentAgreement,
		BusinessCondition:   req.BusinessCondition,
		BeginTime:           time.Unix(req.BeginTime, 0),
		CreatedId:           uid,
		UpdatedId:           uid,
	})

	for i, _ := range projectNodes {
		projectNodes[i].ProjectId = projectId
	}

	for i, _ := range projectPersons {
		projectPersons[i].ProjectId = projectId
	}

	for i, _ := range projectContacts {
		projectContacts[i].ProjectId = projectId
	}

	// todo
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}
	return &RespCreateProject{}, nil
}

func (b *bisLogic) buildProjectNode(finishedNodeIds []int, uid int) ([]*model.ProjectNode, error) {
	finishedNIdMap := make(map[int]struct{}, len(finishedNodeIds))
	for _, v := range finishedNodeIds {
		finishedNIdMap[v] = struct{}{}
	}

	nodes, err := node.BisLogic.Options()
	if err != nil {
		return nil, err
	}
	projectNodes := make([]*model.ProjectNode, 0)
	for _, pNode := range nodes {
		for _, nameItem := range pNode.Children {
			state := model.ProjectNodeStateWaitBegin
			if _, ok := finishedNIdMap[nameItem.Id]; ok {
				state = model.ProjectNodeStateFinished
			}
			projectNodes = append(projectNodes, &model.ProjectNode{
				PId:       pNode.Id,
				NodeId:    nameItem.Id,
				Name:      nameItem.Name,
				IsLeaf:    false,
				State:     state,
				CreatedId: uid,
				UpdatedId: uid,
			})
		}
	}
	return projectNodes, nil
}

func (b *bisLogic) buildProjectPerson(persons []*CreateProjectPerson, uid int) ([]*model.ProjectPerson, error) {
	projectPersons := make([]*model.ProjectPerson, 0)

	ids := make([]int, 0)
	for _, person := range persons {
		ids = append(ids, person.UserId)
	}

	users, err := model.UserModel.SelectByIds(ids)
	if err != nil {
		return nil, err
	}

	userInfoMap := make(map[int]*model.User, len(users))
	for _, user := range users {
		userInfoMap[user.Id] = user
	}

	for _, person := range persons {
		projectPersons = append(projectPersons, &model.ProjectPerson{
			UserId:      person.UserId,
			Type:        person.Type,
			Name:        userInfoMap[person.UserId].Name,
			PhoneNumber: userInfoMap[person.UserId].PhoneNumber,
			CreatedId:   uid,
			UpdatedId:   uid,
		})
	}

	return projectPersons, nil
}

func (b *bisLogic) buildProjectContact(contacts []*CreateProjectContact, uid int) ([]*model.ProjectContact, error) {
	projectContacts := make([]*model.ProjectContact, 0)

	for _, contact := range contacts {
		projectContacts = append(projectContacts, &model.ProjectContact{
			Type:        contact.Type,
			Name:        contact.Name,
			PhoneNumber: contact.PhoneNumber,
			Description: contact.Description,
			CreatedId:   uid,
		})
	}

	return projectContacts, nil
}

func (b *bisLogic) Delete(req *ReqDeleteProject) (*RespDeleteProject, error) {
	if err := model.ProjectModel.Delete(req.Id); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除项目失败")
	}
	return &RespDeleteProject{}, nil
}

func (b *bisLogic) Update(req *ReqUpdateProject, uid int) (*RespUpdateProject, error) {
	err := model.ProjectModel.Update(&model.Project{
		Id:                  req.Id,
		Name:                req.Name,
		Description:         req.Description,
		Attr:                req.Attr,
		State:               0,
		Type:                req.Type,
		Capacity:            req.Capacity,
		Properties:          req.Properties,
		Area:                req.Area,
		Address:             req.Address,
		Connect:             req.Connect,
		InvestmentAgreement: req.InvestmentAgreement,
		BusinessCondition:   req.BusinessCondition,
		BeginTime:           time.Unix(req.BeginTime, 0),
		UpdatedId:           uid,
	})

	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改项目失败")
	}

	return &RespUpdateProject{}, nil
}

func (b *bisLogic) Info(req *ReqInfoProject) (*RespInfoProject, error) {
	project, err := model.ProjectModel.Find(req.Id)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目信息错误")
	}

	res := &RespInfoProject{
		Id:                  project.Id,
		Name:                project.Name,
		Description:         project.Description,
		Attr:                project.Attr,
		State:               project.State,
		Type:                project.Type,
		NodeId:              project.NodeId,
		NodeName:            project.NodeName,
		Schedule:            project.Schedule,
		Capacity:            project.Capacity,
		Properties:          project.Properties,
		Area:                project.Area,
		Address:             project.Address,
		Connect:             project.Connect,
		InvestmentAgreement: project.InvestmentAgreement,
		BusinessCondition:   project.BusinessCondition,
		BeginTime:           project.BeginTime,
		CreatedId:           project.CreatedId,
		UpdatedId:           project.UpdatedId,
		CreatedAt:           project.CreatedAt.Format(global.TimeFormat),
		UpdatedAt:           project.UpdatedAt.Format(global.TimeFormat),
	}

	return res, nil
}

func (b *bisLogic) List(req *ReqProjectList) (*RespProjectList, error) {
	logic.VerifyReqPage(&req.ReqPage)
	cond := b.buildSearchCond(req)
	totalCount, err := model.ProjectModel.GetProjectCountByCond(cond)
	if err != nil {
		slog.Error("list project get user count error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目列表错误")
	}
	projects, err := model.ProjectModel.GetProjectByCond(cond, req.Page, req.Size)
	if err != nil {
		slog.Error("list project get user list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目列表错误")
	}

	items := make([]*ListProjectItem, 0, len(projects))
	for _, pro := range projects {
		items = append(items, &ListProjectItem{
			Id:                  pro.Id,
			Name:                pro.Name,
			Attr:                pro.Attr,
			State:               pro.State,
			Type:                pro.Type,
			NodeName:            pro.NodeName,
			Schedule:            pro.Schedule,
			Capacity:            pro.Capacity,
			Properties:          pro.Properties,
			Area:                pro.Area,
			Address:             pro.Address,
			Connect:             pro.Connect,
			InvestmentAgreement: pro.InvestmentAgreement,
			BusinessCondition:   pro.BusinessCondition,
			BeginTime:           pro.BeginTime.Format(global.TimeFormat),
			CreatedAt:           pro.CreatedAt.Format(global.TimeFormat),
			UpdatedAt:           pro.CreatedAt.Format(global.TimeFormat),
		})
	}

	return &RespProjectList{
		Page:  req.Page,
		Size:  req.Size,
		Count: totalCount,
		List:  items,
	}, nil
}

func (b *bisLogic) buildSearchCond(req *ReqProjectList) *model.ProjectCond {
	cond := &model.ProjectCond{}
	if req.Id > 0 {
		cond.Id = req.Id
	}

	if req.Name != "" {
		cond.Name = req.Name
	}

	if req.NodeId > 0 {
		cond.NodeId = req.NodeId
	}

	if _, ok := model.ProjectStateDesc[req.State]; ok {
		cond.State = req.State
	}

	if _, ok := model.ProjectAttrDesc[req.Attr]; ok {
		cond.Attr = req.Attr
	}

	if _, ok := model.ProjectTypeDesc[req.Type]; ok {
		cond.Type = req.Type
	}

	return cond
}
