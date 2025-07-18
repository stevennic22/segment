package sg

import (
	"log"
	"net/http"
	"segment/core"
	"segment/wsSrv"
	"strings"

	"github.com/gin-gonic/gin"
)

func ConfigureRoutes(rG *gin.RouterGroup, hub *wsSrv.Hub) {
	// Index page, redirects
	rG.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/events/list")
	})

	// Lists recent events
	rG.GET("/list", func(c *gin.Context) {
		// Simply returns a list of the in-memory events from earlier
		c.JSON(http.StatusOK, &struct {
			Events []core.Event `json:"events"`
		}{core.Events})
	})

	// Set up a route to receive events from Segment
	rG.POST("/create", func(c *gin.Context) {
		var event core.Event

		if err := c.ShouldBindJSON(&event); err != nil {
			log.Printf("Binding error: %s", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// Process the received event (in this case we just log it)
		log.Printf("Received event: %+v\n", event)

		// Save the event to our stubbed out database connection
		saveErr := SaveEvent(event)
		if core.CheckError(saveErr, false) {
			if strings.Contains(saveErr.Error(), "Message with Id ") && strings.Contains(saveErr.Error(), ") already exists") {
				c.JSON(http.StatusOK, gin.H{
					"message":       "Duplicate event received.",
					"clients_count": hub.GetClientCount(),
					"event_type":    event.Type,
				})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to process event"})
			return
		}

		// Broadcast to all WebSocket connections
		hub.BroadcastEvent(event)

		c.JSON(http.StatusOK, gin.H{
			"message":       "Event created and broadcast",
			"clients_count": hub.GetClientCount(),
			"event_type":    event.Type,
		})
	})
}
