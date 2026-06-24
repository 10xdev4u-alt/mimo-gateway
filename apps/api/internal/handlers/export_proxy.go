package handlers

import (
	"encoding/csv"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleExportProxyJSON(c *gin.Context) {
	data := []gin.H{
		{"id": "1", "model": "mimo-auto", "status": "active"},
		{"id": "2", "model": "mimo-auto", "status": "active"},
	}

	output, _ := json.MarshalIndent(data, "", "  ")

	c.Header("Content-Type", "application/json")
	c.Header("Content-Disposition", "attachment; filename=export.json")
	c.Data(http.StatusOK, "application/json", output)
}

func HandleExportProxyCSV(c *gin.Context) {
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=export.csv")

	w := csv.NewWriter(c.Writer)
	w.Write([]string{"id", "model", "status"})
	w.Write([]string{"1", "mimo-auto", "active"})
	w.Write([]string{"2", "mimo-auto", "active"})
	w.Flush()
}
