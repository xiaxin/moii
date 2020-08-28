package request

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



type Config struct {
	Generator func() string
}

// New initializes the RequestID middleware.
func New(config ...Config) gin.HandlerFunc {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.Generator == nil {
		cfg.Generator = func() string {
			return uuid.New().String()
		}
	}

	return func(c *gin.Context) {
		rid := c.GetHeader(HeaderXRequestID)
		if rid == "" {
			rid = cfg.Generator()
		}

		// Set the id to ensure that the requestid is in the response
		c.Header(HeaderXRequestID, rid)
		c.Next()
	}
}

// Get returns the request identifier
func Get(c *gin.Context) string {
	return c.Writer.Header().Get(HeaderXRequestID)
}
