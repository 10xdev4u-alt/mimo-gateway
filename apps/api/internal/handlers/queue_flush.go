package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleQueueFlush(c *gin.Context) {
	processor.mu.Lock()
	flushed := len(processor.items)
	processor.items = []interface{}{}
	processor.mu.Unlock()

	c.JSON(http.StatusOK, gin.H{
		"flushed": flushed,
		"status":  "ok",
	})
}
