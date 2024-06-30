package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJSON(c *gin.Context, payload any) error {
	if c.Request == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(c.Request.Body).Decode(payload)
}

func WriteJSON(c *gin.Context, status int, v any) error {
	c.Header("Content-type", "application/json")
	c.Status(status)

	return json.NewEncoder(c.Writer).Encode(v)
}

func WriteError(c *gin.Context, status int, err error) {
	WriteJSON(c, status, gin.H{
		"error": err.Error(),
	})
}
