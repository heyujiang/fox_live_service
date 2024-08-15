package project

import (
	"errors"
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic"
	"fox_live_service/internal/app/server/logic/node"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"strconv"
	"strings"
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

	RespCreateProject struct {
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
		Star                int     `json:"star"`
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
		Id                  int     `json:"id"`
		Name                string  `json:"name"`
		Description         string  `json:"description"`
		Attr                int     `json:"attr"`
		State               int     `json:"state"`
		Type                int     `json:"type"`
		NodeId              int     `json:"nodeId"`
		NodeName            string  `json:"nodName"`
		Schedule            float64 `json:"schedule"`
		Capacity            float64 `json:"capacity"`
		Properties          string  `json:"properties"`
		Area                float64 `json:"area"`
		Address             string  `json:"address"`
		Connect             string  `json:"connect"`
		UserId              int     `json:"userId"`
		Username            string  `json:"username"`
		Star                int     `json:"star"`
		InvestmentAgreement string  `json:"investmentAgreement"`
		BusinessCondition   string  `json:"businessCondition"`
		BeginTime           string  `json:"beginTime"`
		CreatedId           int     `json:"createdId"`
		UpdatedId           int     `json:"updatedId"`
		CreatedAt           string  `json:"createdAt"`
		UpdatedAt           string  `json:"updatedAt"`
	}

	ReqProjectList struct {
		logic.ReqPage
		ReqFromProjectList
	}

	ReqFromProjectList struct {
		Name      string  `form:"name"`
		UserId    int     `form:"userId"`
		Star      int     `form:"star"`
		Type      int     `form:"type"`
		CreatedAt []int64 `form:"createdAt[]"`
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
		NodeName            string  `json:"nodeName"`
		Schedule            float64 `json:"schedule"`
		Capacity            float64 `json:"capacity"`
		Properties          string  `json:"properties"`
		Area                float64 `json:"area"`
		Description         string  `json:"description"`
		Address             string  `json:"address"`
		Connect             string  `json:"connect"`
		Star                int     `json:"star"`
		Username            string  `json:"username"`
		InvestmentAgreement string  `json:"investmentAgreement"`
		BusinessCondition   string  `json:"businessCondition"`
		BeginTime           string  `json:"beginTime"`
		CreatedAt           string  `json:"createdAt"`
		UpdatedAt           string  `json:"updatedAt"`
	}

	ReqProjectOption struct {
	}

	RespProjectOption struct {
		Id   int    `json:"value"`
		Name string `json:"label"`
	}

	RespProjectViewCount struct {
		TotalCount        int     `json:"totalCount"`
		TotalCapacity     float64 `json:"totalCapacity"`
		MonthAddCount     int     `json:"monthAddCount"`
		ThreeStartProject int     `json:"threeStartProject"`
	}

	ReqGetLatestProject struct {
		Type   string `form:"type"`
		Latest string `form:"latest"`
	}

	RespPersonCapacityItem struct {
		UserId   int     `json:"userId"`
		Username string  `json:"username"`
		Capacity float64 `json:"capacity"`
	}
)

func newBisLogic() *bisLogic {
	return &bisLogic{}
}

func (b *bisLogic) Create(req *ReqCreateProject, uid int) (*RespCreateProject, error) {
	if len(req.Person) == 0 {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "项目负责人不能为空")
	}
	firstUserid, firstUsername := 0, ""
	for _, person := range req.Person {
		if person.Type == model.ProjectPersonTypeFirst {
			user, err := model.UserModel.Find(person.UserId)
			if err != nil {
				return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
			}
			firstUserid = user.Id
			firstUsername = user.Username
		}
	}
	if firstUserid == 0 {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "不能没有第一负责人")
	}

	projectNodes, nowNodeId, nowNodeName, schedule, err := b.buildProjectNode(req.NodeIds, uid)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}

	projectContacts, err := b.buildProjectContact(req.Contact, uid)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}

	projectPersons, err := b.buildProjectPerson(req.Person, uid)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}

	projectIns := &model.Project{
		Name:                req.Name,
		Description:         req.Description,
		Attr:                req.Attr,
		State:               req.State,
		Type:                req.Type,
		Capacity:            req.Capacity,
		Schedule:            schedule,
		Properties:          req.Properties,
		Area:                req.Area,
		Address:             req.Address,
		Connect:             req.Connect,
		Star:                req.Star,
		NodeId:              nowNodeId,
		NodeName:            nowNodeName,
		UserId:              firstUserid,
		Username:            firstUsername,
		InvestmentAgreement: req.InvestmentAgreement,
		BusinessCondition:   req.BusinessCondition,
		BeginTime:           time.Unix(req.BeginTime, 0),
		CreatedId:           uid,
		UpdatedId:           uid,
	}
	projectId, err := model.ProjectModel.Create(projectIns)

	for i, _ := range projectNodes {
		projectNodes[i].ProjectId = projectId
	}

	for i, _ := range projectPersons {
		projectPersons[i].ProjectId = projectId
	}

	for i, _ := range projectContacts {
		projectContacts[i].ProjectId = projectId
	}

	if err := model.ProjectContactModel.BatchInsert(projectContacts); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}

	if err := model.ProjectPersonModel.BatchInsert(projectPersons); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}

	if err := model.ProjectNodeModel.BatchInsert(projectNodes); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}

	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建项目失败")
	}
	return &RespCreateProject{
		Id:                  projectId,
		Name:                projectIns.Name,
		Description:         projectIns.Description,
		Attr:                projectIns.Attr,
		Type:                projectIns.Type,
		State:               projectIns.State,
		Capacity:            projectIns.Capacity,
		Properties:          projectIns.Properties,
		Area:                projectIns.Area,
		Star:                projectIns.Star,
		Address:             projectIns.Address,
		Connect:             projectIns.Connect,
		InvestmentAgreement: projectIns.InvestmentAgreement,
		BusinessCondition:   projectIns.BusinessCondition,
		BeginTime:           projectIns.BeginTime.Format(global.DateFormat),
	}, nil
}

func (b *bisLogic) buildProjectNode(finishedNodeIds []int, uid int) ([]*model.ProjectNode, int, string, float64, error) {
	finishedNIdMap := make(map[int]struct{}, len(finishedNodeIds))
	for _, v := range finishedNodeIds {
		finishedNIdMap[v] = struct{}{}
	}

	nodes, err := node.BisLogic.GetAllTreeNodes()
	if err != nil {
		slog.Error("create project build project node error", "err", err)
		return nil, 0, "", 0, errorx.NewErrorX(errorx.ErrCommon, "构建项目节点错误")
	}
	nowNodeId := 0
	nowNodeName := ""
	pNodeCountMap := make(map[int]int, len(nodes))

	projectNodes := make([]*model.ProjectNode, 0)
	var sonNodeCount float64 = 0
	var finishedNodeCount float64 = 0
	for _, n := range nodes {
		for _, v := range n.Children {
			state := model.ProjectNodeStateWaitBegin
			if _, ok := finishedNIdMap[v.Id]; ok {
				state = model.ProjectNodeStateFinished
				pNodeCountMap[n.Id]++
				nowNodeId = v.Id
				nowNodeName = fmt.Sprintf("%s/%s", n.Name, v.Name)
				finishedNodeCount++
			}

			projectNodes = append(projectNodes, &model.ProjectNode{
				PId:       v.Pid,
				NodeId:    v.Id,
				Name:      v.Name,
				IsLeaf:    v.IsLeaf,
				Sort:      v.Sort,
				State:     state,
				CreatedId: uid,
				UpdatedId: uid,
			})
			sonNodeCount++
		}
	}

	for _, v := range nodes {
		state := model.ProjectNodeStateWaitBegin
		sonCount := pNodeCountMap[v.Id]
		if sonCount == len(v.Children) {
			state = model.ProjectNodeStateFinished
		} else if sonCount > 0 {
			state = model.ProjectNodeStateInProcess
		}
		projectNodes = append(projectNodes, &model.ProjectNode{
			PId:       v.Pid,
			NodeId:    v.Id,
			Name:      v.Name,
			IsLeaf:    v.IsLeaf,
			Sort:      v.Sort,
			State:     state,
			CreatedId: uid,
			UpdatedId: uid,
		})
	}
	var schedule float64 = 0
	if sonNodeCount > 0 && finishedNodeCount > 0 {
		schedule, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", finishedNodeCount/sonNodeCount), 64)
	}

	slog.Info("项目进度：", "son node count", sonNodeCount, "finish node count", finishedNodeCount, "schedule", schedule)

	return projectNodes, nowNodeId, nowNodeName, schedule, nil
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
			Name:        userInfoMap[person.UserId].Username,
			PhoneNumber: userInfoMap[person.UserId].PhoneNumber,
			CreatedId:   uid,
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

func (b *bisLogic) Delete(req *ReqDeleteProject, uid int) (*RespDeleteProject, error) {
	_, err := model.ProjectModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目出错")
	}

	if err := model.ProjectModel.Delete(req.Id, uid); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除项目失败")
	}
	return &RespDeleteProject{}, nil
}

func (b *bisLogic) Update(req *ReqUpdateProject, uid int) (*RespUpdateProject, error) {
	_, err := model.ProjectModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "删除项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目出错")
	}

	err = model.ProjectModel.Update(&model.Project{
		Id:                  req.Id,
		Name:                req.Name,
		Description:         req.Description,
		Attr:                req.Attr,
		Star:                req.Star,
		State:               req.State,
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
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目不存在")
		}
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
		UserId:              project.UserId,
		Username:            project.Username,
		Star:                project.Star,
		InvestmentAgreement: project.InvestmentAgreement,
		BusinessCondition:   project.BusinessCondition,
		BeginTime:           project.BeginTime.Format(global.DateFormat),
		CreatedId:           project.CreatedId,
		UpdatedId:           project.UpdatedId,
		CreatedAt:           project.CreatedAt.Format(global.TimeFormat),
		UpdatedAt:           project.UpdatedAt.Format(global.TimeFormat),
	}

	return res, nil
}

func (b *bisLogic) List(req *ReqProjectList, uid int) (*RespProjectList, error) {
	logic.VerifyReqPage(&req.ReqPage)

	projectIds, err := b.getMyProjectIds(uid)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, err.Error())
	}

	cond := b.buildSearchCond(req, projectIds)
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
			Description:         pro.Description,
			Area:                pro.Area,
			Address:             pro.Address,
			Connect:             pro.Connect,
			Star:                pro.Star,
			Username:            pro.Username,
			InvestmentAgreement: pro.InvestmentAgreement,
			BusinessCondition:   pro.BusinessCondition,
			BeginTime:           pro.BeginTime.Format(global.DateFormat),
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

func (b *bisLogic) buildSearchCond(req *ReqProjectList, projectIds []int) *model.ProjectCond {
	cond := &model.ProjectCond{}

	cond.Ids = projectIds

	if req.Name != "" {
		cond.Name = req.Name
	}

	if req.UserId != 0 {
		cond.UserId = req.UserId
	}

	if req.Star != 0 {
		cond.Star = req.Star
	}

	if req.Type != 0 {
		cond.Type = req.Type
	}

	if len(req.CreatedAt) == 2 {
		cond.CreatedAt = []time.Time{time.Unix(req.CreatedAt[0], 0), time.Unix(req.CreatedAt[1], 0)}
	}

	return cond
}

func (b *bisLogic) Option(req *ReqProjectOption) ([]*RespProjectOption, error) {
	projects, err := model.ProjectModel.Select()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目错误")
	}

	res := make([]*RespProjectOption, 0, len(projects))
	for _, v := range projects {
		res = append(res, &RespProjectOption{
			Id:   v.Id,
			Name: v.Name,
		})
	}

	return res, nil
}

func (b *bisLogic) getMyProjectIds(uid int) ([]int, error) {
	//查询用户的项目
	user, err := model.UserModel.Find(uid)
	if err != nil {
		slog.Error("get user info error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取用户信息错误")
	}

	roleIds := cast.ToIntSlice(strings.Split(user.RoleIds, ","))
	for _, roleId := range roleIds {
		if roleId == model.SuperManagerRoleId {
			return []int{}, nil
		}
	}

	// 判断用户角色
	projectPerson, err := model.ProjectPersonModel.SelectByUserId(uid)
	if err != nil {
		slog.Error("get user has project list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取项目列表错误")
	}
	projectIds := make([]int, 0, len(projectPerson))
	for _, person := range projectPerson {
		projectIds = append(projectIds, person.ProjectId)
	}
	return projectIds, nil
}

func (b *bisLogic) GetMyProject(uid int) ([]*ListProjectItem, error) {
	projectIds, err := b.getMyProjectIds(uid)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, err.Error())
	}
	if len(projectIds) == 0 {
		return []*ListProjectItem{}, nil
	}

	projects, err := model.ProjectModel.SelectByIds(projectIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询我的项目出错")
	}

	res := make([]*ListProjectItem, 0, len(projects))
	for _, v := range projects {
		res = append(res, &ListProjectItem{
			Id:                  v.Id,
			Name:                v.Name,
			Attr:                v.Attr,
			State:               v.State,
			Type:                v.Type,
			NodeName:            v.NodeName,
			Schedule:            v.Schedule,
			Capacity:            v.Capacity,
			Properties:          v.Properties,
			Area:                v.Area,
			Description:         v.Description,
			Address:             v.Address,
			Connect:             v.Connect,
			Star:                v.Star,
			Username:            v.Username,
			InvestmentAgreement: v.InvestmentAgreement,
			BusinessCondition:   v.BusinessCondition,
			BeginTime:           v.BeginTime.Format(global.DateFormat),
			CreatedAt:           v.CreatedAt.Format(global.TimeFormat),
			UpdatedAt:           v.UpdatedAt.Format(global.TimeFormat),
		})
	}
	return res, nil
}

func (b *bisLogic) ViewCount() (*RespProjectViewCount, error) {
	projects, err := model.ProjectModel.Select()
	if err != nil {
		slog.Error("get project list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目数据出错")
	}

	res := &RespProjectViewCount{
		TotalCount:        len(projects),
		TotalCapacity:     0.00,
		MonthAddCount:     0,
		ThreeStartProject: 0,
	}

	now := time.Now()
	monthTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	fmt.Println(fmt.Sprintf("month begin : %v", monthTime))
	for _, v := range projects {
		res.TotalCapacity += v.Capacity
		if v.Star == 3 {
			res.ThreeStartProject++
		}

		if v.CreatedAt.After(monthTime) {
			res.MonthAddCount++
		}
	}

	return res, nil
}

func (b *bisLogic) GetLatestProject(req *ReqGetLatestProject) ([]*ListProjectItem, error) {
	beginTime := time.Now()
	if req.Latest == "month" {
		beginTime = beginTime.AddDate(0, -1, 0)
	} else {
		beginTime = beginTime.AddDate(0, 0, -7)
	}
	fmt.Println(fmt.Sprintf("begintime : %+v ; type : %s ", beginTime, req.Latest))
	var notExistIds, existIds, star = make([]int, 0), make([]int, 0), 0

	if req.Type == "threeStar" {
		star = 3
	} else {
		projectIds, err := model.ProjectRecordModel.SelectProjectIdFromCreatedAt(&beginTime)
		if err != nil {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "数据查询出错")
		}
		if req.Type == "noProgress" {
			notExistIds = projectIds
		} else {
			existIds = projectIds
		}
		beginTime = time.Unix(0, 0)
	}

	projects, err := model.ProjectModel.SelectLatestProject(&beginTime, notExistIds, existIds, star)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "数据查询出错")
	}

	res := make([]*ListProjectItem, 0, len(projects))
	for _, v := range projects {
		res = append(res, &ListProjectItem{
			Id:                  v.Id,
			Name:                v.Name,
			Attr:                v.Attr,
			State:               v.State,
			Type:                v.Type,
			NodeName:            v.NodeName,
			Schedule:            v.Schedule,
			Capacity:            v.Capacity,
			Properties:          v.Properties,
			Area:                v.Area,
			Description:         v.Description,
			Address:             v.Address,
			Connect:             v.Connect,
			Star:                v.Star,
			Username:            v.Username,
			InvestmentAgreement: v.InvestmentAgreement,
			BusinessCondition:   v.BusinessCondition,
			BeginTime:           v.BeginTime.Format(global.DateFormat),
			CreatedAt:           v.CreatedAt.Format(global.TimeFormat),
			UpdatedAt:           v.UpdatedAt.Format(global.TimeFormat),
		})
	}
	return res, nil
}

func (b *bisLogic) PersonCapacity() ([]*RespPersonCapacityItem, error) {
	projects, err := model.ProjectModel.Select()
	if err != nil {
		slog.Error("get project list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获得用户容量数据出错")
	}

	projectCap := make(map[int]float64, len(projects))
	projectIds := make([]int, 0)
	for _, v := range projects {
		projectCap[v.Id] += v.Capacity
		projectIds = append(projectIds, v.Id)
	}

	fmt.Println(fmt.Sprintf("%+v", projectCap))

	projectPersons, err := model.ProjectPersonModel.SelectByProjectIds(projectIds)

	resMap := make(map[int]*RespPersonCapacityItem)
	for _, v := range projectPersons {
		fmt.Println(fmt.Sprintf("username : %s , project_id : %d", v.Name, v.ProjectId))
		if _, ok := resMap[v.UserId]; !ok {
			resMap[v.UserId] = &RespPersonCapacityItem{
				UserId:   v.UserId,
				Username: v.Name,
				Capacity: 0,
			}
		}

		resMap[v.UserId].Capacity += projectCap[v.ProjectId]
	}
	res := make([]*RespPersonCapacityItem, 0, len(resMap))

	for _, v := range resMap {
		res = append(res, v)
	}

	return res, nil
}

func (b *bisLogic) ProjectData() error {
	//查询所有项目

	//查询所有员工

	return nil
}
