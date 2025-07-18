package sg

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"segment/core"

	"github.com/gin-gonic/gin"
)

// Save/Log new events
func SaveEvent(event core.Event) error {
	// To Do: Connect DB functions to log event
	// For now, it'll simply use an in-memory slice
	if core.FindEventIndexByMessageID(event.MessageID) != -1 {
		res := fmt.Sprintf("Message with Id (%s) already exists", event.MessageID)
		return fmt.Errorf("%s", res)
	}

	log.Printf("Saving event: %+v\n", event)
	core.Events = append(core.Events, event)

	return nil
}

// Verify X-Signature header
// Untested
func verifySignature(c *gin.Context, secrets *core.Settings) bool {
	signature := c.Request.Header.Get("X-Signature")

	log.Printf("Received Sig: %s", signature)

	bodyBytes, err := io.ReadAll(c.Request.Body)
	c.Request.Body.Close()
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Println(err)
		return false
	}
	digest := []byte(fmt.Sprintf("%s", bodyBytes))
	h := hmac.New(sha256.New, secrets.Secret)
	h.Write(digest)
	signatureDigest := hex.EncodeToString(h.Sum(nil))

	log.Printf("Generated Sig: %s", signatureDigest)

	return signature == string(signatureDigest)
}

// Signature middleware configuration
// Untested
func SigMiddleware(c *gin.Context, secrets *core.Settings) {
	if !verifySignature(c, secrets) {
		c.Set("status", http.StatusUnauthorized)
		c.Set("authorized", false)
	} else {
		c.Set("status", http.StatusOK)
		c.Set("authorized", true)
	}

	log.Printf("Request authorization: %t", c.GetBool("authorized"))
}
