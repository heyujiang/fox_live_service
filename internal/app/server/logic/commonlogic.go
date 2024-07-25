package logic

type (
	ReqPage struct {
		Size int `form:"size" json:"size"`
		Page int `form:"page" json:"page"`
	}
)

func VerifyReqPage(req *ReqPage) {
	if req.Size <= 0 {
		req.Size = 20
	}

	if req.Page <= 0 {
		req.Page = 1
	}
}
