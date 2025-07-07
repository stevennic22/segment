package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"segment/core"
	"segment/database"
	"segment/sg"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var base_file_path string = "/usr/src/app/content/"

var static_file_path string = path.Join(base_file_path, "static")

func generateCorsConfig(h *database.DBHandler) gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOriginFunc = h.AllowedOriginCheck
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	return cors.New(corsConfig)
}

func generateLoggerConfig() gin.HandlerFunc {
	return gin.LoggerWithConfig(
		gin.LoggerConfig{
			SkipPaths: []string{"/state", "/uptime"},
			// Formatter layout
			// [GIN] 2023/04/22 - 22:05:31 | 200 |      62.107µs |     10.10.10.X | GET      "/"
			Formatter: func(param gin.LogFormatterParams) string {
				return fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-8s %#v\n%s",
					param.TimeStamp.Format("2006/01/02 - 15:04:05"),
					param.StatusCode,
					param.Latency,
					param.ClientIP,
					param.Method,
					param.Request.Host+param.Path,
					param.ErrorMessage,
				)
			},
		},
	)
}

func init_logging(environment string) {
	// Configure paths
	logFileDir := path.Join(base_file_path, "logs")
	logFilePath := path.Join(logFileDir, "segment-"+environment+".log")

	// Create directory
	logDirErr := os.MkdirAll(logFileDir, 0775)

	if logDirErr != nil {
		log.Printf("Error opening/creating log directory: %s", logDirErr)
		os.Exit(1)
	}

	// Create file
	f, logFileErr := os.OpenFile(logFilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if logFileErr != nil {
		log.Printf("Error opening/creating log file: %s", logFileErr)
		os.Exit(1)
	}

	// Logs are written to Stdout and to logfile
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
}

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	core.CurrentEnv.Location = path.Join(static_file_path, ".env")

	var environment string = core.CurrentEnv.RetrieveValue("environment")
	init_logging(environment)
}

func main() {
	// Load webhook secret
	secretsConfig := sg.Settings{
		Secret: []byte(core.CurrentEnv.RetrieveValue("webhooksecret")),
	}

	var enableDBVar = core.CurrentEnv.RetrieveValue("db_enable")
	var enableDB bool

	if strings.ToLower(enableDBVar) == "true" {
		enableDB = true
	} else {
		enableDB = false
	}

	var h database.DBHandler

	if enableDB {
		// Connect to Database and add to handler struct
		db, dbInitErr := database.ConnectDB()
		core.CheckError(dbInitErr, true)
		db.SetMaxIdleConns(2)
		db.SetConnMaxLifetime(time.Hour)

		h = *database.NewDBHandler(db)
	}

	// Initialize gin engine
	r := gin.New()

	// Add upstream proxy config
	r.SetTrustedProxies([]string{core.CurrentEnv.RetrieveValue("trusted_proxy")})
	r.TrustedPlatform = "X-Client-IP"

	// Use middlewares
	// Conditionally disable CORS config if no DB
	if enableDB {
		r.Use(
			generateCorsConfig(&h),
			generateLoggerConfig(),
			gin.Recovery(),
		)
	} else {
		r.Use(
			generateLoggerConfig(),
			gin.Recovery(),
		)
	}

	// Index page
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	// For healthcheck
	r.GET("/state", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Segment Event routes
	sg.ConfigureRoutes(r.Group("/events", func(c *gin.Context) { sg.SigMiddleware(c, &secretsConfig) }))

	// Configure port and run
	var serverPort string = ":" + core.CurrentEnv.RetrieveValue("listen_port")
	log.Printf("Starting server at %s\n", serverPort[1:])
	r.Run(serverPort) // listen and serve on 0.0.0.0:serverPort
}
