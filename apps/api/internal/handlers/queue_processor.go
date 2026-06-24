package handlers

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

type QueueProcessor struct {
	mu       sync.Mutex
	items    []interface{}
	processed int
}

var processor = &QueueProcessor{}

func (qp *QueueProcessor) Process(item interface{}) {
	qp.mu.Lock()
	defer qp.mu.Unlock()
	qp.items = append(qp.items, item)
	qp.processed++
}

func HandleQueueProcess(c *gin.Context) {
	var req struct {
		Data interface{} `json:"data"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error())
		return
	}

	processor.Process(req.Data)

	c.JSON(http.StatusOK, gin.H{
		"status":    "processed",
		"total":     processor.processed,
	})
}
