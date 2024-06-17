package services

import "github.com/gin-gonic/gin"

var (
	MessageSuccess         = gin.H{"message": "Success"}
	MessageTimeout         = gin.H{"error": ErrTimeout.Error()}
	MessageInvalidInput    = gin.H{"error": ErrInvalidInput.Error()}
	MessageEmpty           = gin.H{"message": "Empty"}
	MessageInvalidResponse = gin.H{"error": ErrInvalidResponse.Error()}
)
