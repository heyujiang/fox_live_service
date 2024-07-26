package project

var RecordLogic = newRecordLogic()

type (
	recordLogic struct{}

	ReqCreateProjectRecord struct{}

	RespCreateProjectRecord struct{}

	ReqDeleteProjectRecord struct {
		Id int `uri:"id"`
	}

	RespDeleteProjectRecord struct{}

	ReqUpdateProjectRecord struct {
		ReqUriUpdateProjectRecord
		ReqBodyUpdateProjectRecord
	}

	ReqUriUpdateProjectRecord struct {
		Id int `uri:"id"`
	}

	ReqBodyUpdateProjectRecord struct {
	}

	RespUpdateProjectRecord struct {
	}

	ReqInfoProjectRecord struct {
		Id int `uri:"id"`
	}

	RespInfoProjectRecord struct {
	}

	ReqProjectRecordList struct {
	}

	RespProjectRecordList struct {
		List []*ListProjectRecordItem
	}

	ListProjectRecordItem struct {
	}
)

func newRecordLogic() *recordLogic {
	return &recordLogic{}
}

func (b *recordLogic) Create(req *ReqCreateProjectRecord) (*RespCreateProjectRecord, error) {

	return &RespCreateProjectRecord{}, nil
}

func (b *recordLogic) Delete(req *ReqDeleteProjectRecord) (*RespDeleteProjectRecord, error) {
	return &RespDeleteProjectRecord{}, nil
}

func (b *recordLogic) Update(req *ReqUpdateProjectRecord) (*RespUpdateProjectRecord, error) {
	return &RespUpdateProjectRecord{}, nil
}

func (b *recordLogic) Info(req *ReqInfoProjectRecord) (*RespInfoProjectRecord, error) {
	return &RespInfoProjectRecord{}, nil
}

func (b *recordLogic) List(req *ReqProjectRecordList) (*RespProjectRecordList, error) {
	return &RespProjectRecordList{}, nil
}
