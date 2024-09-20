package logic

import (
	"errors"
	"fmt"
	"fox_live_service/internal/app/server/model"
	"github.com/spf13/cast"
	"golang.org/x/exp/slog"
	"strings"
)

type (
	ReqPage struct {
		Size int `form:"size" json:"size"`
		Page int `form:"page" json:"page"`
	}

	RespOption struct {
		Label string `json:"label"`
		Value int    `json:"value"`
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

// CalcProjectProgress 计算获取项目进度
func CalcProjectProgress(projectId int) (float64, error) {
	_, err := model.ProjectModel.Find(projectId)
	if err != nil {
		slog.Error("calc project progress error", "projectId", projectId, "err", err.Error())
		if errors.Is(err, model.ErrNotRecord) {
			return 0, nil
		}
		return 0, err
	}

	nodes, err := model.ProjectNodeModel.GetAllChild(projectId)
	if err != nil {
		slog.Error("calc project progress error", "projectId", projectId, "err", err.Error())
		return 0, err
	}

	var hasProgressTotal float64 = 0
	for _, v := range nodes {
		if v.State == model.ProjectNodeStateInProcess {
			hasProgressTotal += 0.5
		} else if v.State == model.ProjectNodeStateFinished {
			hasProgressTotal += 1
		}
	}
	return cast.ToFloat64(fmt.Sprintf("%.2f", hasProgressTotal/float64(len(nodes)))), nil
}
