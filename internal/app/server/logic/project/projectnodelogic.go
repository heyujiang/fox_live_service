package project

import (
	"errors"
	"fox_live_service/internal/app/server/model"
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
	}

	RespProjectNodeList struct {
		List []*ListProjectNodeItem
	}

	ListProjectNodeItem struct {
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
)

func newNodeLogic() *nodeLogic {
	return &nodeLogic{}
}

func (b *nodeLogic) Create(req *ReqCreateProjectNode) (*RespCreateProjectNode, error) {

	return &RespCreateProjectNode{}, nil
}

func (b *nodeLogic) Delete(req *ReqDeleteProjectNode) (*RespDeleteProjectNode, error) {
	return &RespDeleteProjectNode{}, nil
}

func (b *nodeLogic) Update(req *ReqUpdateProjectNode) (*RespUpdateProjectNode, error) {
	return &RespUpdateProjectNode{}, nil
}

func (b *nodeLogic) Info(req *ReqInfoProjectNode) (*RespInfoProjectNode, error) {
	return &RespInfoProjectNode{}, nil
}

func (b *nodeLogic) List(req *ReqProjectNodeList) (*RespProjectNodeList, error) {
	return &RespProjectNodeList{}, nil
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
