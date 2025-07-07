package sg

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(rG *gin.RouterGroup) {
	// Index page, redirects
	rG.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/events/list")
	})

	// Lists recent events
	rG.GET("/list", func(c *gin.Context) {
		// Simplly returns a list of the in-memory events from earlier
		c.JSON(http.StatusOK, &struct {
			Events []Event `json:"events"`
		}{Events})
	})

	// Set up a route to receive events from Segment
	rG.POST("/create", func(c *gin.Context) {
		var event Event
		if err := c.BindJSON(&event); err != nil {
			log.Printf("Binding error: %s", err)
			c.JSON(400, gin.H{"error": "Invalid request body"})
			return
		}

		// Process the received event (in this case we just log it)
		log.Printf("Received event: %+v\n", event)

		// Save the event to our stubbed out database connection
		if err := SaveEvent(event); err != nil {
			c.JSON(500, gin.H{"error": "Failed to process event"})
			return
		}

		c.Status(200)
	})
}
