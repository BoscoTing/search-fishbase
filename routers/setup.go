package routers

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	UploadDir   = "./uploads"
	ResultsDir  = "./results"
	ServerPort  = ":8080"
	MaxFileSize = 10 << 20 // 10 MB
)

func SetupRouter() *gin.Engine {
	// Create upload and results directory if it doesn't exist
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}
	if err := os.MkdirAll(ResultsDir, 0755); err != nil {
		log.Fatalf("Failed to create results directory: %v", err)
	}

	r := gin.Default() // Initialize router

	r.Static("/results", ResultsDir) // Configure static file serving

	r.LoadHTMLGlob("templates/*")

	registerRoutes(r)

	return r
}
