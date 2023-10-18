package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

func DefaultHeaders() gin.HandlerFunc {

	epoch := time.Unix(0, 0).Format(time.RFC1123)

	defaultHeaders := map[string]string{
		"Expires":                epoch,
		"Cache-Control":          "no-cache, no-store, must-revalidate;",
		"Pragma":                 "no-cache",
		"X-Content-Type-Options": "nosniff",
	}

	// ETag headers array.
	etagHeaders := [6]string{
		"ETag",
		"If-Modified-Since",
		"If-Match",
		"If-None-Match",
		"If-Range",
		"If-Unmodified-Since",
	}

	return func(c *gin.Context) {
		// Delete any ETag headers that may have been set
		for _, v := range etagHeaders {
			if c.Request.Header.Get(v) != "" {
				c.Request.Header.Del(v)
			}
		}

		// Set our Default headers
		for k, v := range defaultHeaders {
			c.Writer.Header().Set(k, v)
		}

		c.Next()

		if c.Writer.Written() {
			c.Writer.Header().Set("Content-Type", "application/json")
		}
	}
}
