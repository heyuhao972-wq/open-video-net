package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type authRequest struct {
	PublicKey string `json:"public_key"`
	Nonce     string `json:"nonce"`
	Signature string `json:"signature"`
}

type registerRequest struct {
	PublicKey string `json:"public_key"`
}

type challengeRequest struct {
	PublicKey string `json:"public_key"`
}

func Register(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}

	var req registerRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	user, err := userService.Register(req.PublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func LoginChallenge(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}

	var req challengeRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	nonce, err := userService.CreateChallenge(req.PublicKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nonce": nonce,
	})
}

func Login(c *gin.Context) {
	if userService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "user service not initialized",
		})
		return
	}

	var req authRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid data",
		})
		return
	}

	token, user, err := userService.Login(req.PublicKey, req.Nonce, req.Signature)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
