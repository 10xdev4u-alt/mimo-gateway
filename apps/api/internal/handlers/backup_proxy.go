package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ProxyBackup struct {
	ID        string      `json:"id"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"created_at"`
}

var proxyBackups = []ProxyBackup{}

func HandleCreateProxyBackup(c *gin.Context) {
	var data interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		BadRequest(c, err.Error())
		return
	}

	backup := ProxyBackup{
		ID:        "backup-" + time.Now().Format("20060102150405"),
		Data:      data,
		CreatedAt: time.Now(),
	}
	proxyBackups = append(proxyBackups, backup)

	c.JSON(http.StatusCreated, gin.H{
		"id":     backup.ID,
		"status": "created",
	})
}

func HandleListProxyBackups(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"backups": proxyBackups,
		"count":   len(proxyBackups),
	})
}
