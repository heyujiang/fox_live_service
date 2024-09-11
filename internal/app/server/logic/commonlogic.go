package logic

import (
	"fmt"
	"strings"
)

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

func GenSortMap(sortStr string) string {
	sortStr = strings.Trim(sortStr, " ")
	if sortStr == "" {
		return ""
	}
	sorts := strings.Split(strings.Trim(sortStr, " "), ",")
	if len(sorts) == 0 {
		return ""
	}
	var sortCond = ""
	for _, s := range sorts {
		fieldSort := strings.Split(s, "_")
		if len(fieldSort) == 2 {
			if fieldSort[1] == "desc" {
				sortCond += fmt.Sprintf(" %s DESC,", fieldSort[0])
			} else {
				sortCond += fmt.Sprintf(" %s ASC,", fieldSort[0])
			}
		} else {
			sortCond += fmt.Sprintf(" %s ASC,", fieldSort[0])
		}
	}

	return sortCond
}
