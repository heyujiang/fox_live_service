package project

var PersonLogic = newPersonLogic()

type (
	personLogic struct{}

	ReqCreateProjectPerson struct{}

	RespCreateProjectPerson struct{}

	ReqDeleteProjectPerson struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectPerson struct{}

	ReqUpdateProjectPerson struct {
		ReqUriUpdateProjectPerson
		ReqBodyUpdateProjectPerson
	}

	ReqUriUpdateProjectPerson struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectPerson struct {
	}

	RespUpdateProjectPerson struct {
	}

	ReqInfoProjectPerson struct {
		Id int `uri:"id"`
	}

	RespInfoProjectPerson struct {
	}

	ReqProjectPersonList struct {
	}

	RespProjectPersonList struct {
		List []*ListProjectPersonItem
	}

	ListProjectPersonItem struct {
	}
)

func newPersonLogic() *personLogic {
	return &personLogic{}
}

func (b *personLogic) Create(req *ReqCreateProjectPerson) (*RespCreateProjectPerson, error) {

	return &RespCreateProjectPerson{}, nil
}

func (b *personLogic) Delete(req *ReqDeleteProjectPerson) (*RespDeleteProjectPerson, error) {
	return &RespDeleteProjectPerson{}, nil
}

func (b *personLogic) Update(req *ReqUpdateProjectPerson) (*RespUpdateProjectPerson, error) {
	return &RespUpdateProjectPerson{}, nil
}

func (b *personLogic) Info(req *ReqInfoProjectPerson) (*RespInfoProjectPerson, error) {
	return &RespInfoProjectPerson{}, nil
}

func (b *personLogic) List(req *ReqProjectPersonList) (*RespProjectPersonList, error) {
	return &RespProjectPersonList{}, nil
}
