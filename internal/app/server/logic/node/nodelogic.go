package node

import (
	"errors"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/logic"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
)

var BisLogic = newBisLogic()

func newBisLogic() *bisLogic {
	return &bisLogic{}
}

type (
	bisLogic struct{}

	RespNodeItem struct {
		Id        int             `json:"id"`
		Name      string          `json:"name"`
		Pid       int             `json:"pid"`
		IsLeaf    int             `json:"isLeaf"`
		Sort      int             `json:"sort"`
		CreatedAt string          `json:"createdAt"`
		UpdatedAT string          `json:"updatedAt"`
		Children  []*RespNodeItem `json:"children"`
	}

	ReqCreateNode struct {
		Name       string `json:"name"`
		Pid        int    `json:"pid"`
		Sort       int    `json:"sort"`
		SyncAll    int    `json:"syncAll"`
		ProjectIds []int  `json:"projectIds"`
	}

	RespCreateNode struct{}

	ReqDeleteNode struct {
		Id int `uri:"id"`
	}

	RespDeleteNode struct{}

	ReqUpdateNode struct {
		ReqUriUpdateNode
		ReqBodyUpdateNode
	}

	ReqUriUpdateNode struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateNode struct {
		Name string `json:"name"`
		Sort int    `json:"sort"`
	}

	RespUpdateNode struct {
	}

	ReqNodeInfo struct {
		Id int `uri:"id"`
	}

	RespNodeInfo struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Pid       int    `json:"pid"`
		IsLeaf    int    `json:"isLeaf"`
		Sort      int    `json:"sort"`
		CreatedAt string `json:"createAt"`
		UpdatedAT string `json:"updateAt"`
	}

	RespParentNodeItem struct {
		Id       int                   `json:"id"`
		Name     string                `json:"name"`
		Pid      int                   `json:"pid"`
		Sort     int                   `json:"sort"`
		Children []*RespParentNodeItem `json:"children"`
	}

	RespOptionItem struct {
		Id   int    `json:"value"`
		Name string `json:"label"`
		//Disabled bool              `json:"disabled"`
		Children []*RespOptionItem `json:"children,omitempty"`
	}

	RespParentNodeOptionItem struct {
		Id   int    `json:"value"`
		Name string `json:"label"`
	}

	Item struct {
		Id        int    `json:"id"`
		Name      string `json:"name"`
		Pid       int    `json:"pid"`
		IsLeaf    int    `json:"isLeaf"`
		Sort      int    `json:"sort"`
		CreatedId int    `json:"createdId"`
		UpdatedId int    `json:"updatedId"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}

	TreeItem struct {
		*Item
		Children []*TreeItem `json:"children"`
	}
)

// Create 创建节点
func (b *bisLogic) Create(req *ReqCreateNode, uid int) (*RespCreateNode, error) {
	isLeaf := model.NodeLeafNo
	if req.Pid > 0 {
		isLeaf = model.NodeLeafYes

		_, err := model.NodeModel.Find(req.Pid)
		if err != nil {
			slog.Error("create user get user error ", "id", req.Pid, "err", err)
			if errors.Is(err, model.ErrNotRecord) {
				return nil, errorx.NewErrorX(errorx.ErrCommon, "父节点不存在")
			}
			return nil, errorx.NewErrorX(errorx.ErrCommon, "查询父节点错误")
		}
	}
	insertNode := model.Node{
		Name:      req.Name,
		Pid:       req.Pid,
		IsLeaf:    isLeaf,
		Sort:      req.Sort,
		CreatedId: uid,
		UpdatedId: uid,
	}

	if err := model.NodeModel.Insert(&insertNode); err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "创建节点失败")
	}

	projectIds := req.ProjectIds
	if req.SyncAll == 1 {
		projects, err := model.ProjectModel.GetAllNoFinished()
		if err != nil {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "获取未完成项目出错")
		}
		projectIds = make([]int, 0, len(projectIds))
		for _, v := range projects {
			projectIds = append(projectIds, v.Id)
		}
	}
	if len(projectIds) > 0 {
		// 同步项目节点
		for _, v := range projectIds {
			go b.syncProjectNode(&insertNode, v, uid)
		}
	}

	return &RespCreateNode{}, nil
}

func (b *bisLogic) syncProjectNode(node *model.Node, projectId int, uid int) {
	if node.Pid > 0 {
		pNode, err := model.ProjectNodeModel.FindByProjectIdAndNodeId(projectId, node.Pid)
		if err != nil {
			slog.Error("sync project node info error ", "err", err.Error())
		}

		if pNode.State == model.ProjectNodeStateFinished { //修改状态为进行中
			if err := model.ProjectNodeModel.UpdateProjectNodeState(pNode.Id, model.ProjectNodeStateInProcess, uid); err != nil {
				slog.Error("sync project node info error ", "err", err.Error())
			}
		}
	}

	if err := model.ProjectNodeModel.Insert(&model.ProjectNode{
		ProjectId: projectId,
		PId:       node.Pid,
		NodeId:    node.Id,
		Name:      node.Name,
		IsLeaf:    node.IsLeaf,
		Sort:      node.Sort,
		State:     model.ProjectNodeStateWaitBegin,
		CreatedId: uid,
		UpdatedId: uid,
	}); err != nil {
		slog.Error("sync project node info error ", "err", err.Error())
	}

	schedule, err := logic.CalcProjectProgress(projectId)
	if err != nil {
		slog.Error("sync project node info error ", "err", err.Error())
	}

	if err := model.ProjectModel.UpdateSchedule(projectId, schedule, uid); err != nil {
		slog.Error("sync project node info error ", "err", err.Error())
	}
}

// Delete 删除节点
func (b *bisLogic) Delete(req *ReqDeleteNode, uid int) (*RespDeleteNode, error) {
	node, err := model.NodeModel.Find(req.Id)
	if err != nil {
		slog.Error("delete node get node error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "节点不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点错误")
	}

	if node.IsLeaf == model.NodeLeafNo {
		nodes, err := model.NodeModel.SelectByPid(node.Id)
		if err != nil {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "查询子节点错误")
		}
		if len(nodes) > 0 {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "存在子节点，不可删除")
		}
	}

	if err := model.NodeModel.Delete(req.Id); err != nil {
		slog.Error("delete node error ", "id", req.Id, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除节点信息错误")
	}

	projects, err := model.ProjectModel.GetAllNoFinished()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取未完成项目出错")
	}

	for _, v := range projects {
		go b.delProjectNode(v.Id, req.Id, uid)
	}

	return &RespDeleteNode{}, nil
}

func (b *bisLogic) delProjectNode(projectId, nodeId, uid int) {
	if err := model.ProjectNodeModel.DelNode(projectId, nodeId); err != nil {
		slog.Error("delete project node error ", "id", nodeId, "err", err)
	}

	if err := model.ProjectRecordModel.DeleteByProjectIdAndNodeId(projectId, nodeId, uid); err != nil {
		slog.Error("delete project node error ", "id", nodeId, "err", err)
	}

	if err := model.ProjectAttachedModel.DeleteByProjectIdAndNodeId(projectId, nodeId); err != nil {
		slog.Error("delete project node error ", "id", nodeId, "err", err)
	}

	schedule, err := logic.CalcProjectProgress(projectId)
	if err != nil {
		slog.Error("sync project node info error ", "err", err.Error())
	}

	if err := model.ProjectModel.UpdateSchedule(projectId, schedule, uid); err != nil {
		slog.Error("sync project node info error ", "err", err.Error())
	}
}

// Update 修改节点信息
func (b *bisLogic) Update(req *ReqUpdateNode, uid int) (*RespUpdateNode, error) {
	_, err := model.NodeModel.Find(req.Id)
	if err != nil {
		slog.Error("update node get node error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "节点不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点错误")
	}

	if err := model.NodeModel.Update(&model.Node{
		Id:        req.Id,
		Name:      req.Name,
		Sort:      req.Sort,
		UpdatedId: uid,
	}); err != nil {
		slog.Error("update node error ", "id", req.Id, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "修改节点信息错误")
	}

	return &RespUpdateNode{}, nil
}

// Info 节点信息
func (b *bisLogic) Info(req *ReqNodeInfo) (*RespNodeInfo, error) {
	node, err := model.NodeModel.Find(req.Id)
	if err != nil {
		slog.Error("get user error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "节点不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点错误")
	}

	return &RespNodeInfo{
		Id:        node.Id,
		Name:      node.Name,
		Pid:       node.Pid,
		IsLeaf:    node.IsLeaf,
		CreatedAt: node.CreatedAt.Format(global.TimeFormat),
		UpdatedAT: node.UpdatedAt.Format(global.TimeFormat),
	}, nil
}

func (b *bisLogic) List() ([]*RespNodeItem, error) {
	nodes, err := b.GetAllTreeNodes()
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}

	res := make([]*RespNodeItem, 0)

	for _, node := range nodes {
		children := make([]*RespNodeItem, 0, len(node.Children))
		for _, child := range node.Children {
			children = append(children, &RespNodeItem{
				Id:        child.Id,
				Name:      child.Name,
				Pid:       child.Pid,
				IsLeaf:    child.IsLeaf,
				Sort:      child.Sort,
				CreatedAt: child.CreatedAt,
				UpdatedAT: child.UpdatedAt,
				Children:  make([]*RespNodeItem, 0),
			})
		}
		res = append(res, &RespNodeItem{
			Id:        node.Id,
			Name:      node.Name,
			Pid:       node.Pid,
			IsLeaf:    node.IsLeaf,
			Sort:      node.Sort,
			CreatedAt: node.CreatedAt,
			UpdatedAT: node.UpdatedAt,
			Children:  children,
		})
	}
	return res, nil
}

func (b *bisLogic) Parent() ([]*RespParentNodeOptionItem, error) {
	nodes, err := model.NodeModel.SelectNotLeaf()
	if err != nil {
		slog.Error("parent node get node list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取节点列表错误")
	}

	res := make([]*RespParentNodeOptionItem, 0)
	for _, node := range nodes {
		res = append(res, &RespParentNodeOptionItem{
			Id:   node.Id,
			Name: node.Name,
		})
	}

	return res, nil
}

func (b *bisLogic) Options() ([]*RespOptionItem, error) {
	nodes, err := model.NodeModel.Select()
	if err != nil {
		slog.Error("list node get node list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取节点列表错误")
	}

	nodeIdMap := make(map[int][]*RespOptionItem)
	nodeIdMap[0] = make([]*RespOptionItem, 0)
	for _, node := range nodes {
		if _, ok := nodeIdMap[node.Pid]; !ok {
			nodeIdMap[node.Pid] = make([]*RespOptionItem, 0)
		}

		nodeIdMap[node.Pid] = append(nodeIdMap[node.Pid], &RespOptionItem{
			Id:   node.Id,
			Name: node.Name,
		})
	}
	for _, nodeList := range nodeIdMap {
		for _, node := range nodeList {
			if _, ok := nodeIdMap[node.Id]; ok {
				node.Children = nodeIdMap[node.Id]
			}
		}
	}

	return nodeIdMap[0], nil
}

// GetAllNodes 获取所有节点数据
func (b *bisLogic) GetAllNodes() ([]*Item, error) {
	nodes, err := model.NodeModel.Select()
	if err != nil {
		slog.Error("list node get node list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}

	nodeItems := make([]*Item, 0, len(nodes))
	for _, node := range nodes {
		nodeItems = append(nodeItems, &Item{
			Id:        node.Id,
			Name:      node.Name,
			Pid:       node.Pid,
			IsLeaf:    node.IsLeaf,
			Sort:      node.Sort,
			CreatedId: node.CreatedId,
			UpdatedId: node.UpdatedId,
			CreatedAt: node.CreatedAt.Format(global.TimeFormat),
			UpdatedAt: node.UpdatedAt.Format(global.TimeFormat),
		})
	}

	return nodeItems, nil
}

func (b *bisLogic) GetAllTreeNodes() ([]*TreeItem, error) {
	nodes, err := b.GetAllNodes()
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}

	nodePIdMap := make(map[int][]*TreeItem)
	for _, node := range nodes {
		if _, ok := nodePIdMap[node.Pid]; !ok {
			nodePIdMap[node.Pid] = make([]*TreeItem, 0)
		}
		nodePIdMap[node.Pid] = append(nodePIdMap[node.Pid], &TreeItem{
			Item:     node,
			Children: make([]*TreeItem, 0),
		})
	}

	for _, nodeList := range nodePIdMap {
		for _, node := range nodeList {
			if _, ok := nodePIdMap[node.Id]; ok {
				node.Children = nodePIdMap[node.Id]
			}
		}
	}

	return nodePIdMap[0], nil
}

func (b *bisLogic) ProjectOptions() ([]*logic.RespOption, error) {
	projects, err := model.ProjectModel.GetAllNoFinished()
	if err != nil {
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目列表出错")
	}

	res := make([]*logic.RespOption, 0, len(projects))
	for _, project := range projects {
		res = append(res, &logic.RespOption{
			Label: project.Name,
			Value: project.Id,
		})
	}

	return res, nil
}
