package project

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
