package upload

import (
	"fmt"
	"fox_live_service/config/global"
	"fox_live_service/internal/app/server/model"
	"fox_live_service/pkg/errorx"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/exp/slog"
	"mime/multipart"
	"path"
	"time"
)

var (
	BisLogic = bisLogic{}
)

type (
	bisLogic struct{}

	ReqFileUpload struct {
		File *multipart.FileHeader `form:"file" binding:"required"`
		Type string                `form:"type" binding:"required"`
	}

	RespFileUpload struct {
		Url      string `json:"url"`
		Filename string `json:"filename"`
		Size     int64  `json:"size"`
		Mime     string `json:"mime"`
	}

	FileFilter struct {
		MaxSize int64
		Mime    map[string]struct{}
	}
)

const (
	FileUploadUserAvatar         = "avatar"
	FileUploadRecordAttachedFile = "attached"
)

var (
	Mime2Suffix = map[string]string{
		"image/jpeg":                  ".jpg",
		"image/png":                   ".png",
		"image/gif":                   ".gif",
		"image/bmp":                   ".bmp",
		"image/webp":                  ".webp",
		"application/pdf":             "pdf",
		"application/zip":             "zip",
		"application/vnd.rar":         "rar",
		"application/x-7z-compressed": "7z",
		"application/vnd.ms-excel":    "xls",
		"application/wps-office.xls":  "xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "xlsx",
		"application/wps-office.xlsx":                                               "xlsx",
		"application/msword":                                                        "doc",
		"application/wps-office.doc":                                                "docx",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "docx",
		"application/wps-office.docx":                                               "docx",
		"application/vnd.ms-powerpoint":                                             "ppt",
		"application/wps-office.ppt":                                                "ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
		"application/wps-office.pptx":                                               "pptx",
	}

	FileFilerMap = map[string]*FileFilter{
		FileUploadUserAvatar: {
			MaxSize: 10 * 1000 * 1000,
			Mime: map[string]struct{}{
				"image/jpeg": {},
				"image/png":  {},
				"image/gif":  {},
				"image/bmp":  {},
				"image/webp": {},
			},
		},
		FileUploadRecordAttachedFile: {
			MaxSize: 100 * 1000 * 1000,
			Mime: map[string]struct{}{
				"image/jpeg":                  {},
				"image/png":                   {},
				"image/gif":                   {},
				"image/bmp":                   {},
				"image/webp":                  {},
				"application/pdf":             {},
				"application/zip":             {},
				"application/vnd.rar":         {},
				"application/x-7z-compressed": {},
				"application/vnd.ms-excel":    {},
				"application/wps-office.xls":  {},
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         {},
				"application/wps-office.xlsx":                                               {},
				"application/msword":                                                        {},
				"application/wps-office.doc":                                                {},
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   {},
				"application/wps-office.docx":                                               {},
				"application/vnd.ms-powerpoint":                                             {},
				"application/wps-office.ppt":                                                {},
				"application/vnd.openxmlformats-officedocument.presentationml.presentation": {},
				"application/wps-office.pptx":                                               {},
			},
		},
	}

	FileUploadPath = map[string]string{
		FileUploadRecordAttachedFile: "attached",
		FileUploadUserAvatar:         "avatar",
	}
)

func newBisLogic() *bisLogic {
	return &bisLogic{}
}

// Upload 文件上传接口
func (b *bisLogic) Upload(c *gin.Context, req *ReqFileUpload, uid int) (*RespFileUpload, error) {
	err := b.filter(req)
	if err != nil {
		return nil, err
	}

	mime := req.File.Header.Get("Content-Type")
	fileSavePath := b.getUploadPath(req.Type, Mime2Suffix[mime])

	err = c.SaveUploadedFile(req.File, fileSavePath)
	if err != nil {
		slog.Error("upload file error", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "上传文件出错")
	}
	url := fmt.Sprintf("http://%s:%d/%s", global.Config.GetString("Host"), global.Config.GetInt("Port"), fileSavePath)
	if err := model.FileModel.Insert(&model.File{
		Type:      req.Type,
		Url:       url,
		Mime:      mime,
		Size:      req.File.Size,
		Filename:  req.File.Filename,
		CreatedId: uid,
	}); err != nil {
		slog.Error("upload file error", "err", err)
		return nil, errorx.NewErrorX(errorx.ErrCommon, "上传文件保存文件信息失败")
	}

	return &RespFileUpload{
		Url:      url,
		Filename: req.File.Filename,
		Size:     req.File.Size,
		Mime:     req.File.Header.Get("Content-Type"),
	}, nil
}

func (b *bisLogic) filter(req *ReqFileUpload) error {
	filter, ok := FileFilerMap[req.Type]
	if !ok {
		return errorx.NewErrorX(errorx.ErrCommon, "不支持的上传类型")
	}

	if req.File.Size > filter.MaxSize {
		return errorx.NewErrorX(errorx.ErrCommon, "文件大小超出限制")
	}
	fileMime := req.File.Header.Get("Content-Type")
	fmt.Println("file mime is ", fileMime)
	_, ok = filter.Mime[fileMime]
	if !ok {
		return errorx.NewErrorX(errorx.ErrCommon, "不支持的文件格式")
	}
	return nil
}

func (b *bisLogic) getUploadPath(fileType string, ext string) string {
	fileName := fmt.Sprintf("%s.%s", uuid.NewV4().String(), ext)
	return path.Join(global.FileUploadPath, FileUploadPath[fileType], time.Now().Format(global.DateDirFormat), fileName)
}
