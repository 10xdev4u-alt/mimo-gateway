package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func HandleDiffText(c *gin.Context) {
	var req struct {
		Text1 string `json:"text1"`
		Text2 string `json:"text2"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequestError(c, err.Error())
		return
	}

	lines1 := strings.Split(req.Text1, "\n")
	lines2 := strings.Split(req.Text2, "\n")

	added := 0
	removed := 0
	for i, line := range lines1 {
		if i < len(lines2) {
			if line != lines2[i] {
				removed++
				added++
			}
		} else {
			removed++
		}
	}
	removed += len(lines2) - len(lines1)

	c.JSON(http.StatusOK, gin.H{
		"lines1":  len(lines1),
		"lines2":  len(lines2),
		"added":   added,
		"removed": removed,
	})
}
