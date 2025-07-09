package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"segment/core"
	"segment/sg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func AllowedOriginCheck(origin string) bool {
	trusted_domain := core.Config.RetrieveValue("trusted_domain")

	regexStr := fmt.Sprintf("^http(s)*://(.*.)*(%s)$", trusted_domain)
	r := regexp.MustCompile(regexStr)

	if r.MatchString(origin) {
		return true
	} else {
		return false
	}
}

func generateCorsConfig() gin.HandlerFunc {
	corsConfig := cors.DefaultConfig()

	corsConfig.AllowOriginFunc = AllowedOriginCheck
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
	logFileDir := path.Join(core.Config.BaseFilePath, "logs")
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

	core.Config.Location = path.Join(core.Config.FilePath("static"), ".env")
	core.Config.InitEnvFile()

	var environment string = core.Config.RetrieveValue("environment")
	init_logging(environment)
}

func main() {
	// Load webhook secret
	secretsConfig := core.Settings{
		Secret: []byte(core.Config.RetrieveValue("webhooksecret")),
	}

	// Initialize gin engine
	r := gin.New()

	// Add upstream proxy config
	r.SetTrustedProxies([]string{core.Config.RetrieveValue("trusted_proxy")})
	r.TrustedPlatform = "X-Client-IP"

	// Use middlewares
	r.Use(
		generateCorsConfig(),
		generateLoggerConfig(),
		gin.Recovery(),
	)

	// Index page
	r.GET("/", func(c *gin.Context) { c.Status(http.StatusOK) })

	// For healthcheck
	r.GET("/state", func(c *gin.Context) { c.Status(http.StatusOK) })

	r.StaticFile("/favicon.ico", "./favicon.ico")

	r.Static("/content", path.Join(core.Config.BaseFilePath, "direct"))

	// Segment Event routes
	sg.ConfigureRoutes(r.Group("/events", func(c *gin.Context) { sg.SigMiddleware(c, &secretsConfig) }))

	// Configure port and run
	core.Config.FetchServerPort()
	log.Printf("Starting server at %s\n", core.Config.ServerPort[1:])
	r.Run(core.Config.ServerPort) // listen and serve on 0.0.0.0:serverPort
}
