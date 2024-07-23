package logic

type (
	ReqPage struct {
		Size int `form:"size" json:"size"`
		Page int `form:"page" json:"page"`
	}
)
