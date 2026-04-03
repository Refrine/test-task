package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        startTime := time.Now()
        
        c.Next()
        
        latency := time.Since(startTime)
        statusCode := c.Writer.Status()
        
        logrus.WithFields(logrus.Fields{
            "method":     c.Request.Method,
            "path":       c.Request.URL.Path,
            "status":     statusCode,
            "latency_ms": latency.Milliseconds(),
            "client_ip":  c.ClientIP(),
        }).Info("ok")
    }
}