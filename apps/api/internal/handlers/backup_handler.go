package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Backup struct {
	ID        string      `json:"id"`
	Data      interface{} `json:"data"`
	CreatedAt time.Time   `json:"created_at"`
}

var backups = []Backup{}

func HandleCreateBackup(c *gin.Context) {
	var data interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		BadRequestError(c, err.Error())
		return
	}

	backup := Backup{
		ID:        "backup-" + time.Now().Format("20060102150405"),
		Data:      data,
		CreatedAt: time.Now(),
	}
	backups = append(backups, backup)

	c.JSON(http.StatusCreated, gin.H{
		"id":      backup.ID,
		"status":  "created",
	})
}

func HandleListBackups(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"backups": backups,
		"count":   len(backups),
	})
}
