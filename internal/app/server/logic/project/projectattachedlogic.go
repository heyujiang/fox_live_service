package project

var AttachedLogic = newAttachedLogic()

type (
	attachedLogic struct{}

	ReqCreateProjectAttached struct{}

	RespCreateProjectAttached struct{}

	ReqDeleteProjectAttached struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectAttached struct{}

	ReqUpdateProjectAttached struct {
		ReqUriUpdateProjectAttached
		ReqBodyUpdateProjectAttached
	}

	ReqUriUpdateProjectAttached struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectAttached struct {
	}

	RespUpdateProjectAttached struct {
	}

	ReqInfoProjectAttached struct {
		Id int `uri:"id"`
	}

	RespInfoProjectAttached struct {
	}

	ReqProjectAttachedList struct {
	}

	RespProjectAttachedList struct {
		List []*ListProjectAttachedItem
	}

	ListProjectAttachedItem struct {
	}
)

func newAttachedLogic() *attachedLogic {
	return &attachedLogic{}
}

func (b *attachedLogic) Create(req *ReqCreateProjectAttached) (*RespCreateProjectAttached, error) {

	return &RespCreateProjectAttached{}, nil
}

func (b *attachedLogic) Delete(req *ReqDeleteProjectAttached) (*RespDeleteProjectAttached, error) {
	return &RespDeleteProjectAttached{}, nil
}

func (b *attachedLogic) Update(req *ReqUpdateProjectAttached) (*RespUpdateProjectAttached, error) {
	return &RespUpdateProjectAttached{}, nil
}

func (b *attachedLogic) Info(req *ReqInfoProjectAttached) (*RespInfoProjectAttached, error) {
	return &RespInfoProjectAttached{}, nil
}

func (b *attachedLogic) List(req *ReqProjectAttachedList) (*RespProjectAttachedList, error) {
	return &RespProjectAttachedList{}, nil
}
