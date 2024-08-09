package project

import (
	"errors"
	"fmt"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"golang.org/x/exp/slog"
	"strconv"
)

var NodeLogic = newNodeLogic()

type (
	nodeLogic struct{}

	ReqCreateProjectNode struct{}

	RespCreateProjectNode struct{}

	ReqDeleteProjectNode struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectNode struct{}

	ReqUpdateProjectNode struct {
		ReqUriUpdateProjectNode
		ReqBodyUpdateProjectNode
	}

	ReqUriUpdateProjectNode struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectNode struct {
	}

	RespUpdateProjectNode struct {
	}

	ReqInfoProjectNode struct {
		Id int `uri:"id"`
	}

	RespInfoProjectNode struct {
	}

	ReqProjectNodeList struct {
		ProjectId int `uri:"id"`
	}

	RespProjectNodeItem struct {
		Id            int                    `json:"id"`
		NodeId        int                    `json:"nodeId"`
		Name          string                 `json:"name"`
		Pid           int                    `json:"pid"`
		State         int                    `json:"state"`
		Sort          int                    `json:"sort"`
		IsLeaf        int                    `json:"isLeaf"`
		RecordTotal   int                    `json:"recordTotal"`
		AttachedTotal int                    `json:"attachedTotal"`
		Progress      float64                `json:"progress"`
		Children      []*RespProjectNodeItem `json:"children"`
	}

	NodeItem struct {
		Id     int    `json:"id"`
		NodeId int    `json:"node_id"`
		Name   string `json:"name"`
		Pid    int    `json:"pid"`
		State  int    `json:"state"`
		Sort   int    `json:"sort"`
		IsLeaf int    `json:"is_leaf"`
	}

	TreeNodeItem struct {
		*NodeItem
		Children []*TreeNodeItem `json:"children"`
	}

	ReqProjectNodeOption struct {
		ProjectId int `uri:"id"`
	}

	RespProjectNodeOption struct {
		NodeId   int                      `json:"value"`
		Name     string                   `json:"label"`
		Disabled bool                     `json:"disabled"`
		Children []*RespProjectNodeOption `json:"children"`
	}
)

func newNodeLogic() *nodeLogic {
	return &nodeLogic{}
}

func (n *nodeLogic) Create(req *ReqCreateProjectNode) (*RespCreateProjectNode, error) {

	return &RespCreateProjectNode{}, nil
}

func (n *nodeLogic) Delete(req *ReqDeleteProjectNode) (*RespDeleteProjectNode, error) {
	return &RespDeleteProjectNode{}, nil
}

func (n *nodeLogic) Update(req *ReqUpdateProjectNode) (*RespUpdateProjectNode, error) {
	return &RespUpdateProjectNode{}, nil
}

func (n *nodeLogic) Info(req *ReqInfoProjectNode) (*RespInfoProjectNode, error) {
	return &RespInfoProjectNode{}, nil
}

func (n *nodeLogic) List(req *ReqProjectNodeList) ([]*RespProjectNodeItem, error) {
	if req.ProjectId == 0 {
		return []*RespProjectNodeItem{}, nil
	}

	_, err := model.ProjectModel.Find(req.ProjectId)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目出错")
	}

	nodes, err := n.GetAllTreeNodes(req.ProjectId)
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}

	records, err := model.ProjectRecordModel.GetAllByProjectId(req.ProjectId)
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}
	recordCountMap := make(map[int]int)
	for _, v := range records {
		recordCountMap[v.NodeId]++
	}

	attacheds, err := model.ProjectAttachedModel.GetAllByProjectId(req.ProjectId)
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}
	attachedCountMap := make(map[int]int)
	for _, v := range attacheds {
		attachedCountMap[v.NodeId]++
	}

	res := make([]*RespProjectNodeItem, 0)

	for _, node := range nodes {
		children := make([]*RespProjectNodeItem, 0, len(node.Children))
		var progress, nodeFinishedTotal, nodeTotal float64
		for _, child := range node.Children {
			nodeTotal++
			if child.State == model.ProjectNodeStateInProcess {
				nodeFinishedTotal += 0.5
			} else if child.State == model.ProjectNodeStateFinished {
				nodeFinishedTotal += 1
			}
			children = append(children, &RespProjectNodeItem{
				Id:            child.Id,
				NodeId:        child.NodeId,
				Name:          child.Name,
				Pid:           child.Pid,
				State:         child.State,
				IsLeaf:        child.IsLeaf,
				Sort:          child.Sort,
				Progress:      0,
				RecordTotal:   recordCountMap[child.NodeId],
				AttachedTotal: attachedCountMap[child.NodeId],
				Children:      make([]*RespProjectNodeItem, 0),
			})
		}

		progress, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", nodeFinishedTotal/nodeTotal), 64)

		res = append(res, &RespProjectNodeItem{
			Id:            node.Id,
			NodeId:        node.NodeId,
			Name:          node.Name,
			Pid:           node.Pid,
			State:         node.State,
			IsLeaf:        node.IsLeaf,
			Sort:          node.Sort,
			Progress:      progress,
			RecordTotal:   recordCountMap[node.NodeId],
			AttachedTotal: attachedCountMap[node.NodeId],
			Children:      children,
		})
	}
	return res, nil
}

func (n *nodeLogic) GetAllProjectNode(projectId int) ([]*NodeItem, error) {
	projectNodes, err := model.ProjectNodeModel.SelectByProjectId(projectId)
	if err != nil {
		return nil, errors.New("查询节点所有节点错误")
	}

	res := make([]*NodeItem, 0, len(projectNodes))
	for _, v := range projectNodes {
		res = append(res, &NodeItem{
			Id:     v.Id,
			NodeId: v.NodeId,
			Name:   v.Name,
			Pid:    v.PId,
			State:  v.State,
			Sort:   v.Sort,
			IsLeaf: v.IsLeaf,
		})
	}
	return res, nil
}

func (n *nodeLogic) GetAllTreeNodes(projectId int) ([]*TreeNodeItem, error) {
	nodes, err := n.GetAllProjectNode(projectId)
	if err != nil {
		slog.Error("list project node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目节点数据错误")
	}

	nodePIdMap := make(map[int][]*TreeNodeItem)
	for _, node := range nodes {
		if _, ok := nodePIdMap[node.Pid]; !ok {
			nodePIdMap[node.Pid] = make([]*TreeNodeItem, 0)
		}
		nodePIdMap[node.Pid] = append(nodePIdMap[node.Pid], &TreeNodeItem{
			NodeItem: node,
			Children: make([]*TreeNodeItem, 0),
		})
	}

	for _, nodeList := range nodePIdMap {
		for _, node := range nodeList {
			if _, ok := nodePIdMap[node.NodeId]; ok {
				node.Children = nodePIdMap[node.NodeId]
			}
		}
	}

	return nodePIdMap[0], nil
}

func (n *nodeLogic) Option(req *ReqProjectNodeOption) ([]*RespProjectNodeOption, error) {
	if req.ProjectId == 0 {
		return []*RespProjectNodeOption{}, nil
	}

	_, err := model.ProjectModel.Find(req.ProjectId)
	if err != nil {
		if errors.Is(err, model.ErrNotRecord) {
			return nil, errorx.NewErrorX(errorx.ErrCommon, "项目不存在")
		}
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询项目出错")
	}

	nodes, err := n.GetAllTreeNodes(req.ProjectId)
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}

	records, err := model.ProjectRecordModel.GetAllByProjectId(req.ProjectId)
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}
	recordCountMap := make(map[int]int)
	for _, v := range records {
		recordCountMap[v.NodeId]++
	}

	attacheds, err := model.ProjectAttachedModel.GetAllByProjectId(req.ProjectId)
	if err != nil {
		slog.Error("list node error ", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "查询节点数据错误")
	}
	attachedCountMap := make(map[int]int)
	for _, v := range attacheds {
		attachedCountMap[v.NodeId]++
	}

	res := make([]*RespProjectNodeOption, 0)

	for _, node := range nodes {
		children := make([]*RespProjectNodeOption, 0, len(node.Children))
		for _, child := range node.Children {
			children = append(children, &RespProjectNodeOption{
				NodeId:   child.NodeId,
				Name:     child.Name,
				Children: make([]*RespProjectNodeOption, 0),
			})
		}

		res = append(res, &RespProjectNodeOption{
			NodeId:   node.NodeId,
			Name:     node.Name,
			Disabled: true,
			Children: children,
		})
	}
	return res, nil
}
