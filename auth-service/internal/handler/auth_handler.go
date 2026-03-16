package handler

import (
	"github.com/djwhocodes/auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{s}
}

func (h *AuthHandler) Register(c *gin.Context) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	err := h.service.Register(c, req.Email, req.Password)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "user created"})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	access, refresh, err :=
		h.service.Login(c, req.Email, req.Password)

	if err != nil {
		c.JSON(401, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(200, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}
