package home

import (
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"time"
)

var BisLogic = newBisLogic()

type (
	bisLogic struct{}

	RespUserData struct {
		UserNames      []string `json:"usernames"`
		ProjectCounts  []int    `json:"projectCounts"`
		RecordCounts   []int    `json:"recordCounts"`
		AttachedCounts []int    `json:"attachedCounts"`
	}

	ReqUserData struct {
		TimeRange []int64 `form:"timeRange[]"`
	}
)

func newBisLogic() *bisLogic {
	return &bisLogic{}
}

func (b *bisLogic) UserData(req *ReqUserData) (*RespUserData, error) {
	if len(req.TimeRange) != 2 {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "请选择时间范围")
	}
	startTimes := time.UnixMilli(req.TimeRange[0])
	endTime := time.UnixMilli(req.TimeRange[1])

	//查询所有员工
	users, err := model.UserModel.Select()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户出错")
	}

	userIds := make([]int, 0, len(users))
	userNameMap := make(map[int]string, len(users))
	recordMap := make(map[int]int, len(users))
	projectMap := make(map[int]int, len(users))
	attachedMap := make(map[int]int, len(users))

	for _, user := range users {
		userIds = append(userIds, user.Id)
		userNameMap[user.Id] = user.Username
		recordMap[user.Id] = 0
		projectMap[user.Id] = 0
		attachedMap[user.Id] = 0
	}

	//查询所有用户作为第一项目负责人的项目数
	projectItems, err := model.ProjectPersonModel.SelectGroupCountByUserIds(userIds)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户项目数量出错")
	}
	for _, projectItem := range projectItems {
		projectMap[projectItem.UserId] += projectItem.Count
	}

	//查询所有用户的记录数量
	recordItems, err := model.ProjectRecordModel.SelectGroupCountByUserIds(userIds, startTimes, endTime)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询用户记录数量出错")
	}
	for _, recordItem := range recordItems {
		recordMap[recordItem.UserId] += recordItem.Count
	}

	//查询所有用户的文件数量
	attachedItems, err := model.ProjectAttachedModel.SelectGroupCountByUserIds(userIds, startTimes, endTime)
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询文件项目数量出错")
	}
	for _, attachedItem := range attachedItems {
		attachedMap[attachedItem.UserId] += attachedItem.Count
	}

	res := &RespUserData{
		UserNames:      make([]string, 0),
		ProjectCounts:  make([]int, 0),
		RecordCounts:   make([]int, 0),
		AttachedCounts: make([]int, 0),
	}

	for _, userId := range userIds {
		if projectMap[userId] > 0 || recordMap[userId] > 0 || attachedMap[userId] > 0 {
			res.UserNames = append(res.UserNames, userNameMap[userId])
			res.ProjectCounts = append(res.ProjectCounts, projectMap[userId])
			res.RecordCounts = append(res.RecordCounts, recordMap[userId])
			res.AttachedCounts = append(res.AttachedCounts, attachedMap[userId])
		}
	}

	return res, nil
}
