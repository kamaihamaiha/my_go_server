package server

import (
	"github.com/gin-gonic/gin"

	"my_law_server/internal/handler"
)

func NewRouter(lawHandler *handler.LawHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/healthz", lawHandler.Healthz)

	api := router.Group("/api/v1")
	api.GET("/types/previews", lawHandler.ListTypePreviews)
	api.GET("/types/:typeId/laws", lawHandler.ListLawsByType)
	api.GET("/laws/:versionId/parsed", lawHandler.GetParsedLaw)

	return router
}
