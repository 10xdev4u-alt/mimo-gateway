package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func HandleGetUser(c *gin.Context) {
	user := User{
		ID:    "1",
		Name:  "Admin",
		Email: "admin@mimo-gateway.dev",
		Role:  "admin",
	}
	c.JSON(http.StatusOK, user)
}

func HandleUpdateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "updated",
		"name":    req.Name,
		"email":   req.Email,
	})
}
