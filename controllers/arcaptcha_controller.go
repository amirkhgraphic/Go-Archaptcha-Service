package controllers

import (
	"net/http"

	"github.com/amirkhgraphic/go-arcaptcha-service/services"
	"github.com/gin-gonic/gin"
)

// GenerateFakeChallenge provides a throwaway challenge_id for testing local flows.
func GenerateFakeChallenge(c *gin.Context) {
	token := services.Arcaptcha.GenerateChallenge()
	c.JSON(http.StatusOK, gin.H{
		"challenge_id": token,
		"note":         "Use this challenge_id in protected requests. Suffix -neterr to simulate network errors.",
	})
}

// VerifyFakeChallenge lets you check a token without consuming it.
func VerifyFakeChallenge(c *gin.Context) {
	var body struct {
		ChallengeID string `json:"challenge_id"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	err := services.Arcaptcha.PeekChallenge(body.ChallengeID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"valid": false, "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true})
}
