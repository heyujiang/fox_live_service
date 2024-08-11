package system

import (
	"errors"
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
)

var (
	DeptLogic = newDeptLogic()
)

type (
	deptLogic struct {
	}

	ReqCreateDept struct {
		Title  string `json:"title"`
		Pid    int    `json:"pid"`
		Remark string `json:"remark"`
		Order  int    `json:"order"`
	}

	RespCreateDept struct{}

	ReqUpdateDept struct {
		ReqUpdateDeptUri
		ReqUpdateDeptBody
	}

	ReqUpdateDeptUri struct {
		Id int `uri:"id"`
	}

	ReqUpdateDeptBody struct {
		Title  string `json:"title"`
		Pid    int    `json:"pid"`
		Remark string `json:"remark"`
		Order  int    `json:"order"`
	}

	RespUpdateDept struct{}

	ReqDeleteDept struct {
		Id int `uri:"id"`
	}
	RespDeleteDept struct{}

	ReqUpdateDeptStatus struct {
		ReqUpdateDeptUri
		ReqUpdateDeptStatusBody
	}

	ReqUpdateDeptStatusBody struct {
		Status int `json:"status"`
	}

	RespUpdateDeptStatus struct {
	}

	RespDeptListItem struct {
		Id        int                 `json:"id"`
		Title     string              `json:"title"`
		Pid       int                 `json:"pid"`
		Status    int                 `json:"status"`
		Remark    string              `json:"remark"`
		Order     int                 `json:"order"`
		CreatedId int                 `json:"createdId"`
		UpdatedId int                 `json:"updatedId"`
		CreatedAt string              `json:"createdAt"`
		UpdatedAt string              `json:"updatedAt"`
		Children  []*RespDeptListItem `json:"children"`
	}

	RespDeptParentItem struct {
		Id       int                   `json:"id"`
		Title    string                `json:"title"`
		Children []*RespDeptParentItem `json:"children"`
	}

	ReqGetDeptRules struct {
		Id int `uri:"id"`
	}
)

func newDeptLogic() *deptLogic {
	return &deptLogic{}
}

func (r *deptLogic) Create(req *ReqCreateDept, uid int) (*RespCreateDept, error) {
	if err := model.DeptModel.Create(&model.Dept{
		Name:      req.Title,
		Pid:       req.Pid,
		Status:    model.DeptStatusEnable,
		Remark:    req.Remark,
		Order:     req.Order,
		CreatedId: uid,
		UpdatedId: uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建部门出错")
	}
	return &RespCreateDept{}, nil
}

func (r *deptLogic) Update(req *ReqUpdateDept, uid int) (*RespUpdateDept, error) {
	_, err := model.DeptModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "部门不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询部门信息出错")
	}

	if err := model.DeptModel.Update(&model.Dept{
		Id:        req.Id,
		Name:      req.Title,
		Pid:       req.Pid,
		Remark:    req.Remark,
		Order:     req.Order,
		UpdatedId: uid,
	}); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "编辑部门出错")
	}

	return &RespUpdateDept{}, nil
}

func (r *deptLogic) Delete(req *ReqDeleteDept) (*RespDeleteDept, error) {
	_, err := model.DeptModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "部门不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询部门信息出错")
	}

	if err := model.DeptModel.Delete(req.Id); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除部门出错")
	}
	return &RespDeleteDept{}, nil
}

func (r *deptLogic) List() ([]*RespDeptListItem, error) {
	depts, err := model.DeptModel.Select()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询部门出错")
	}

	deptMap := make(map[int][]*RespDeptListItem)
	deptMap[0] = []*RespDeptListItem{}
	for _, dept := range depts {
		fmt.Println(fmt.Sprintf("%+v", dept))
		if _, ok := deptMap[dept.Pid]; !ok {
			deptMap[dept.Pid] = make([]*RespDeptListItem, 0)
		}

		deptMap[dept.Pid] = append(deptMap[dept.Pid], &RespDeptListItem{
			Id:        dept.Id,
			Title:     dept.Name,
			Pid:       dept.Pid,
			Remark:    dept.Remark,
			Status:    dept.Status,
			Order:     dept.Order,
			CreatedId: dept.CreatedId,
			UpdatedId: dept.UpdatedId,
			CreatedAt: dept.CreatedAt.Format(global.TimeFormat),
			UpdatedAt: dept.UpdatedAt.Format(global.TimeFormat),
		})
	}

	for _, depts := range deptMap {
		for _, dept := range depts {
			if _, ok := deptMap[dept.Id]; !ok {
				dept.Children = make([]*RespDeptListItem, 0)
			} else {
				dept.Children = deptMap[dept.Id]
			}
		}
	}

	return deptMap[0], nil
}

func (r *deptLogic) UpdateStatus(req *ReqUpdateDeptStatus, uid int) (*RespUpdateDeptStatus, error) {
	_, err := model.DeptModel.Find(req.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "部门不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询部门信息出错")
	}

	if err := model.DeptModel.UpdateStatus(req.Id, req.Status, uid); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改部门状态出错")
	}
	return &RespUpdateDeptStatus{}, nil
}

func (r *deptLogic) Parents() ([]*RespDeptParentItem, error) {
	depts, err := model.DeptModel.SelectEnable()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询菜单出错")
	}

	deptMap := make(map[int][]*RespDeptParentItem)
	deptMap[0] = []*RespDeptParentItem{}
	for _, dept := range depts {
		if _, ok := deptMap[dept.Pid]; !ok {
			deptMap[dept.Pid] = make([]*RespDeptParentItem, 0)
		}
		deptMap[dept.Pid] = append(deptMap[dept.Pid], &RespDeptParentItem{
			Id:    dept.Id,
			Title: dept.Name,
		})
	}

	for _, depts := range deptMap {
		for _, dept := range depts {
			if _, ok := deptMap[dept.Id]; !ok {
				dept.Children = make([]*RespDeptParentItem, 0)
			} else {
				dept.Children = deptMap[dept.Id]
			}
		}
	}

	return deptMap[0], nil
}
