package routers

import (
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	r.GET("/", homeHandler) // Home page

	r.POST("/upload", uploadHandler)

	r.GET("/status/:jobId", statusHandler) // Status check

	r.GET("/download/:filename", downloadHandler)
}
