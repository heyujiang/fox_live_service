package router

import (
	"fox_live_service/config/global"
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerStaticRouters(e *gin.Engine) {
	e.StaticFS("/storages/upload", http.Dir(global.FileUploadPath))
}
