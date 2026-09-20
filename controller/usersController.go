package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/scientist-v08/crackers/service"
)

func SignUp(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	if err := service.SignUp(req.Email, req.Password); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "user already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Creating a new user: Successful"})
}

func AdminSignUp(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		IsAdmin  bool   `json:"isAdmin"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	if !req.IsAdmin {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only admins can use this API"})
		return
	}

	if err := service.AdminSignUp(req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Creating a new user: Successful"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	token, routes, err := service.Login(req.Email, req.Password)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "invalid email ID" || err.Error() == "invalid password" {
			status = http.StatusUnauthorized // more accurate
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": token,
		"routes":       routes,
	})
}