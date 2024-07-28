package node

import (
	"errors"
	"fox_live_service/config/global"
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
		Name string `json:"name"`
		Pid  int    `json:"pid"`
		Sort int    `json:"sort"`
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

	return &RespCreateNode{}, nil
}

// Delete 删除节点
func (b *bisLogic) Delete(req *ReqDeleteNode) (*RespDeleteNode, error) {
	_, err := model.NodeModel.Find(req.Id)
	if err != nil {
		slog.Error("delete node get node error ", "id", req.Id, "err", err)
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "节点不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点错误")
	}

	if err := model.NodeModel.Delete(req.Id); err != nil {
		slog.Error("delete node error ", "id", req.Id, "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "删除节点信息错误")
	}

	return &RespDeleteNode{}, nil
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
	nodes, err := model.NodeModel.Select()
	if err != nil {
		slog.Error("list node get node list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取节点列表错误")
	}

	nodeIdMap := make(map[int][]*RespNodeItem)
	nodeIdMap[0] = make([]*RespNodeItem, 0)
	for _, node := range nodes {
		if _, ok := nodeIdMap[node.Pid]; !ok {
			nodeIdMap[node.Pid] = make([]*RespNodeItem, 0)
		}
		nodeIdMap[node.Pid] = append(nodeIdMap[node.Pid], &RespNodeItem{
			Id:        node.Id,
			Name:      node.Name,
			Pid:       node.Pid,
			IsLeaf:    node.IsLeaf,
			Sort:      node.Sort,
			CreatedAt: node.CreatedAt.Format(global.TimeFormat),
			UpdatedAT: node.UpdatedAt.Format(global.TimeFormat),
		})
	}
	for _, nodeList := range nodeIdMap {
		for _, node := range nodeList {
			if _, ok := nodeIdMap[node.Id]; ok {
				node.Children = nodeIdMap[node.Id]
			} else {
				node.Children = make([]*RespNodeItem, 0)
			}
		}
	}

	return nodeIdMap[0], nil
}

func (b *bisLogic) Parent() ([]*RespParentNodeItem, error) {
	nodes, err := model.NodeModel.SelectNotLeaf()
	if err != nil {
		slog.Error("parent node get node list error", "err", err.Error())
		return nil, errorx.NewErrorX(errorx.ErrCommon, "获取节点列表错误")
	}

	nodeIdMap := make(map[int][]*RespParentNodeItem)
	nodeIdMap[0] = make([]*RespParentNodeItem, 0)
	for _, node := range nodes {
		if _, ok := nodeIdMap[node.Pid]; !ok {
			nodeIdMap[node.Pid] = make([]*RespParentNodeItem, 0)
		}
		nodeIdMap[node.Pid] = append(nodeIdMap[node.Pid], &RespParentNodeItem{
			Id:   node.Id,
			Name: node.Name,
			Pid:  node.Pid,
			Sort: node.Sort,
		})
	}
	for _, nodeList := range nodeIdMap {
		for _, node := range nodeList {
			if _, ok := nodeIdMap[node.Id]; ok {
				node.Children = nodeIdMap[node.Id]
			} else {
				node.Children = make([]*RespParentNodeItem, 0)
			}
		}
	}

	return nodeIdMap[0], nil
}
