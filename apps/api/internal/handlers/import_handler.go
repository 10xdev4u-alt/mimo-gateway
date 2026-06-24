package handlers

import (
	"encoding/csv"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleImportCSV(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		BadRequestError(c, "no file provided")
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		BadRequestError(c, "invalid CSV format")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rows":    len(records),
		"headers": records[0],
		"status":  "imported",
	})
}

func HandleImportJSON(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		BadRequestError(c, "failed to read body")
		return
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		BadRequestError(c, "invalid JSON")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "imported",
		"type":   strings.ReplaceAll(strings.TrimLeft(reflect.TypeOf(data).String(), "*"), ".", "_"),
	})
}
